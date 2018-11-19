package cmd

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func wsBuild(wsaddr, msg string) {

	// fmt.Println(wsaddr)

	c, err := gosocketio.Dial(
		wsaddr,
		transport.GetDefaultWebsocketTransport(),
	)
	defer c.Close()

	if err != nil {
		log.Println(err.Error())
	}

	// wg.Add(1)

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
			// wg.Done()
		}
	})

	if err != nil {
		log.Println(err.Error())
	}

	c.Emit("build", msg)

	// wg.Wait()
	// c.Close()

	time.Sleep(900 * time.Second)

}

func wsGetServers(wsaddr string) {

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

	err = c.On("get_servers", func(h *gosocketio.Channel, args string) {
		log.Println("servers: ", args)

		if strings.ContainsAny(args, "{") {
			c.Close()
		}

		if strings.ContainsAny(args, "[") {
			c.Close()
		}

	})

	c.Emit("get_servers", "")

	time.Sleep(1000 * time.Second)
	c.Close()

}

func wsSSH(wsaddr, msg string) {

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

	err = c.On("exec", func(h *gosocketio.Channel, args string) {
		log.Println("output: ", args)

	})

	c.Emit("exec", msg)

	time.Sleep(1000 * time.Second)
	c.Close()

}

func wsGetNodes(wsaddr string) {

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

	err = c.On("get_nodes", func(h *gosocketio.Channel, args string) {
		log.Println("nodes: ", args)

		if strings.ContainsAny(args, "{") {
			c.Close()
		}

		if strings.ContainsAny(args, "[") {
			c.Close()
		}

	})

	c.Emit("get_nodes", "")

	time.Sleep(1000 * time.Second)
	c.Close()

}

func wsGethCmd(wsaddr, cmd string) {

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

	err = c.On(cmd, func(h *gosocketio.Channel, args string) {
		log.Println("Output: ", args)

		match, _ := regexp.MatchString("[a-zA-Z0-9]+", args)
		if match {
			c.Close()
		}

	})

	c.Emit(cmd, "")

	time.Sleep(1000 * time.Second)
	c.Close()

}
