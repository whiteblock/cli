package cmd

import (
	"log"
	"net/url"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func wsBuild(wsaddr, msg string) {

	u := url.URL{Scheme: "ws", Host: wsaddr, Path: "/"}

	c, err := gosocketio.Dial(
		u.String(),
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
	})
	c.Emit("build", msg)

	time.Sleep(1000 * time.Second)
	c.Close()

}

func wsGetServers(wsaddr string) {

	u := url.URL{Scheme: "ws", Host: wsaddr, Path: "/"}

	c, err := gosocketio.Dial(
		u.String(),
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

	err = c.On("server", func(h *gosocketio.Channel, args string) {
		log.Println("get_servers: ", args)
	})

	c.Emit("build", `{"Servers":[4],"Blockchain":"ethereum","Nodes":3,"Image":"ethereum:latest"}`)

	time.Sleep(1000 * time.Second)
	c.Close()

}
