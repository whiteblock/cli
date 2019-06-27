package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/graarh/golang-socketio"
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
		defer util.DeleteStore(".in_progress_build_id")
		res, err := util.JsonRpcCall("stop_build", []string{testnetId})
		if err != nil {
			util.PrintErrorFatal(err)
		}
		fmt.Printf("\r\n%v\r\n", res)
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
				fmt.Printf("\r\n%v\r\n", res)
				signal.Reset(syscall.SIGTSTP)
				syscall.Kill(syscall.Getpid(), syscall.SIGSTOP)
				signal.Notify(pauseChan, syscall.SIGTSTP)
			} else if sigId == syscall.SIGCONT && paused {
				paused = false
				res, err := util.JsonRpcCall("unfreeze_build", []string{testnetId})
				if err != nil {
					util.PrintErrorFatal(err)
				}
				fmt.Printf("\r\n%v\r\n", res)
			}
		}
	}()

	var mutex = &sync.Mutex{}
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
		util.PrintStringError(args)
		os.Exit(1)
	})

	err = c.On("build_status", func(h *gosocketio.Channel, args string) {
		var status BuildStatus
		err := json.Unmarshal([]byte(args), &status)
		if err != nil {
			util.PrintStringError(args)
			os.Exit(1)
		}
		if status.Frozen {
			fmt.Printf("\nBuild is currently frozen. Press Ctrl-\\ to drop into console. Run 'whiteblock build unfreeze' to resume. \r")
		} else if status.Error != nil {
			fmt.Println() //move to the next line
			what := status.Error["what"]
			util.PrintStringError(what)
			os.Exit(1)
			mutex.Unlock()
		} else if status.Progress == 0.0 {
			fmt.Printf("Sending build context to Whiteblock\r")
		} else if status.Progress == 100.0 {
			fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\r", "Build", status.Progress)
			fmt.Println("\a")
			mutex.Unlock()
		} else {
			fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\r", status.Stage, status.Progress)
		}
	})
	if err != nil {
		util.PrintErrorFatal(err)
	}
	c.Emit("build_status", testnetId)
	mutex.Lock()
}
