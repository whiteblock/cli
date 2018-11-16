package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func wsBuild(wsaddr, msg string) {

	done := false

	c, err := gosocketio.Dial(
		wsaddr,
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Println(err.Error())
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
		log.Println("build_status: ", args)

		if args != "Not Ready" {
			done = true
		}
	})
	c.Emit("build", msg)

	if done == true {
		c.Close()
	}

}

func wsGetServers(wsaddr string) {

	done := false

	c, err := gosocketio.Dial(
		wsaddr,
		transport.GetDefaultWebsocketTransport(),
	)

	fmt.Sprintln(gosocketio.GetUrl("localhost", 5000, false))

	if err != nil {
		panic(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		panic(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		panic(err)
	}

	err = c.On("get_servers", func(h *gosocketio.Channel, args string) {
		log.Println("servers: ", args)

		if strings.ContainsAny(args, "{") {
			done = true
		}
	})
	if err != nil {
		panic(err)
	}

	c.Emit("get_servers", "")

	if done == true {
		c.Close()
	}

}
