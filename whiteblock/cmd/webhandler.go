package cmd

import (
	"encoding/json"
	"github.com/graarh/golang-socketio"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type BuildStatus struct {
	Error    map[string]string `json:"error"`
	Progress float64           `json:"progress"`
	Stage    string            `json:"stage"`
	Frozen   bool              `json:"frozen"`
}

func buildListener(testnetId string) {
	sigChan := make(chan os.Signal, 1)
	pauseChan := make(chan os.Signal, 1)
	quitChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT) //Stop the build on SIGINT
	go func() {
		<-sigChan
		defer util.Delete("in_progress_build_id")
		res, err := util.JsonRpcCall("stop_build", []string{testnetId})
		if err != nil {
			util.PrintErrorFatal(err)
		}
		util.Printf("\r\n%v\r\n", res)
		os.Exit(0)
	}()

	signal.Notify(quitChan, syscall.SIGQUIT) //^\ means exit without side effects
	go func() {
		<-quitChan
		os.Exit(0)
	}()

	signal.Notify(pauseChan, syscall.SIGTSTP, syscall.SIGCONT)
	paused := false
	go func() {
		for {
			sigId := <-pauseChan
			if sigId == syscall.SIGTSTP && !paused {
				paused = true
				res, err := util.JsonRpcCall("freeze_build", []string{testnetId})
				if err != nil {
					util.PrintErrorFatal(err)
				}
				util.Printf("\r\n%v\r\n", res)
				signal.Reset(syscall.SIGTSTP)
				syscall.Kill(syscall.Getpid(), syscall.SIGSTOP)
				signal.Notify(pauseChan, syscall.SIGTSTP)
			} else if sigId == syscall.SIGCONT && paused {
				paused = false
				res, err := util.JsonRpcCall("unfreeze_build", []string{testnetId})
				if err != nil {
					util.PrintErrorFatal(err)
				}
				util.Printf("\r\n%v\r\n", res)
			}
		}
	}()

	mutex := &sync.Mutex{}
	mutex.Lock()
	c, err := func() (*gosocketio.Client, error) {
		if strings.HasSuffix(conf.ServerAddr, "5000") { //5000 is http
			return gosocketio.Dial(
				"ws://"+conf.ServerAddr+"/socket.io/?EIO=3&transport=websocket",
				GetDefaultWebsocketTransport())
		} else { //5001 is https
			return gosocketio.Dial(
				"wss://"+conf.ServerAddr+"/socket.io/?EIO=3&transport=websocket",
				GetDefaultWebsocketTransport())
		}
	}()
	if err != nil {
		util.PrintErrorFatal(err)
	}
	defer c.Close()

	c.On("error", func(h *gosocketio.Channel, args string) {
		util.PrintErrorFatal(args)
	})

	err = c.On("build_status", func(h *gosocketio.Channel, args string) {
		var status build.Status
		err := json.Unmarshal([]byte(args), &status)
		if err != nil {
			util.PrintErrorFatal(args)
		}
		if status.Print() {
			mutex.Unlock()
		}
	})
	if err != nil {
		util.PrintErrorFatal(err)
	}
	c.Emit("build_status", testnetId)
	mutex.Lock()
}
