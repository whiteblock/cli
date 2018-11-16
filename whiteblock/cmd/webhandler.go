package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func wsEmitListen(wsaddr, msg string) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: wsaddr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	fmt.Println(u.String())

	dialer := websocket.Dialer{
		HandshakeTimeout: 60 * time.Second,
		Proxy:            http.ProxyFromEnvironment,
	}

	c, resp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("%d", resp.StatusCode)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			log.Printf("recv: %s", message)
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(msg))
}
