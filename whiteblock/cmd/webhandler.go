package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type BuildStatus struct {
	Error    error   `json:"error"`
	Progress float64 `json:"progress"`
}

func wsEmitListen(wsaddr, cmd, param string) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	c, err := gosocketio.Dial(
		wsaddr,
		transport.GetDefaultWebsocketTransport(),
	)

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	defer c.Close()

	// build
	if cmd == "build" {
		err = c.On("build", func(h *gosocketio.Channel, args string) {
			log.Println("build: ", args)
		})

		err = c.On("build_status", func(h *gosocketio.Channel, args string) {
			var status BuildStatus
			json.Unmarshal([]byte(args), &status)
			fmt.Printf("Building: %f \t\t\t\t\r", status.Progress)
			if status.Progress == 100.0 {
				fmt.Println("\nDone")
				mutex.Unlock()
			} else if status.Error != nil {
				fmt.Println(status.Error.Error())
				mutex.Unlock()
			}
		})
	}

	// get servers
	if cmd == "get_servers" {
		err = c.On("get_servers", func(h *gosocketio.Channel, args string) {
			print(args)
			mutex.Unlock()
		})
	}

	// get nodes
	if cmd == "get_nodes" {
		err = c.On("get_nodes", func(h *gosocketio.Channel, args string) {
			print(args)
			mutex.Unlock()
		})
	}

	// gethcmd
	if strings.HasPrefix(cmd, "eth::") {
		err = c.On(cmd, func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				println(args)
				mutex.Unlock()
			}
		})
	}

	// netconfig
	if cmd == "netconfig" {
		err = c.On(cmd, func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				print(args)
			}
			mutex.Unlock()
		})
	}

	// ssh
	if cmd == "exec" {
		err = c.On("exec", func(h *gosocketio.Channel, args string) {
			print(args)
			mutex.Unlock()
		})
	}

	// stats
	if cmd == "stats" {
		err = c.On("stats", func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				print(prettyp(args))
			} else {
				println(err.Error())
			}
			mutex.Unlock()
		})
	}

	// all_stats
	if cmd == "all_stats" {
		err = c.On("all_stats", func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				print(prettyp(args))
			} else {
				println(err.Error())
			}
			mutex.Unlock()
		})
	}

	c.Emit(cmd, param)
	mutex.Lock()
}
