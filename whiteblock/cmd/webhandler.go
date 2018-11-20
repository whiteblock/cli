package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type BuildStatus struct {
	Error    error   `json:"error"`
	Progress float64 `json:"progress"`
}

func wsBuild(wsaddr, msg string) {

	// fmt.Println(wsaddr)
	var mutex = &sync.Mutex{}
	mutex.Lock()
	c, err := gosocketio.Dial(
		wsaddr,
		transport.GetDefaultWebsocketTransport(),
	)
	defer c.Close()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})

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
		}
	})

	if err != nil {
		log.Println(err.Error())
	}

	c.Emit("build", msg)
	mutex.Lock()

}

func wsGetServers(wsaddr string) {
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

	err = c.On("get_servers", func(h *gosocketio.Channel, args string) {
		print(args)
		mutex.Unlock()
	})

	c.Emit("get_servers", "")
	mutex.Lock()
}

func wsSSH(wsaddr, msg string) {
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
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		//log.Println("Connected")
	})

	err = c.On("exec", func(h *gosocketio.Channel, args string) {
		print(args)
		mutex.Unlock()
	})

	c.Emit("exec", msg)
	mutex.Lock()

}

func wsGetNodes(wsaddr string) {
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
	err = c.On("get_nodes", func(h *gosocketio.Channel, args string) {
		print(args)
		mutex.Unlock()
	})

	c.Emit("get_nodes", "")

	mutex.Lock()
}

func wsSendCmd(wsaddr, cmd, param string) {
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

	err = c.On(cmd, func(h *gosocketio.Channel, args string) {
		if len(args) > 0 {
			println(args)
			mutex.Unlock()
		}
	})

	println(cmd)
	println(param)

	c.Emit(cmd, param)

	mutex.Lock()
}
