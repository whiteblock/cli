package cmd

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func wsEmitListen(wsaddr, msg string) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: wsaddr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(msg))
}

// func wsEmit(wsaddr, msg string) {
// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt)

// 	socket := gowebsocket.New(wsaddr)

// 	socket.OnConnected = func(socket gowebsocket.Socket) {
// 		log.Println("Connected to server")
// 	}

// 	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
// 		log.Println("Recieved connect error ", err)
// 	}

// 	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
// 		log.Println("Recieved message " + message)
// 	}

// 	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
// 		log.Println("Recieved binary data ", data)
// 	}

// 	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
// 		log.Println("Recieved ping " + data)
// 	}

// 	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
// 		log.Println("Recieved pong " + data)
// 	}

// 	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
// 		log.Println("Disconnected from server ")
// 		return
// 	}

// 	socket.Connect()

// 	socket.SendText(msg)

// for {
// 	select {
// 	case <-interrupt:
// 		log.Println("interrupt")
// 		socket.Close()
// 		return
// 	}
// }
// }
