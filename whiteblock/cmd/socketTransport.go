package cmd

/*
(The MIT License)

Copyright (c) 2016 Gennadiy Kovalev <gennadiy@weird.company>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/graarh/golang-socketio/transport"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	upgradeFailed = "Upgrade failed: "
)

var (
	ErrorMethodNotAllowed = errors.New("Method not allowed")
)

type WebsocketConnection struct {
	socket    *websocket.Conn
	transport *WebsocketTransport
}

func (wsc *WebsocketConnection) GetMessage() (message string, err error) {
	wsc.socket.SetReadDeadline(time.Now().Add(wsc.transport.ReceiveTimeout))
	msgType, reader, err := wsc.socket.NextReader()
	if err != nil {
		return "", err
	}

	//support only text messages exchange
	if msgType != websocket.TextMessage {
		return "", errors.New("Binary messages are not supported")
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", errors.New("Buffer error")
	}
	text := string(data)

	//empty messages are not allowed
	if len(text) == 0 {
		return "", errors.New("Wrong packet type error")
	}

	return text, nil
}

func (wsc *WebsocketConnection) WriteMessage(message string) error {
	wsc.socket.SetWriteDeadline(time.Now().Add(wsc.transport.SendTimeout))
	writer, err := wsc.socket.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func (wsc *WebsocketConnection) Close() {
	wsc.socket.Close()
}

func (wsc *WebsocketConnection) PingParams() (interval, timeout time.Duration) {
	return wsc.transport.PingInterval, wsc.transport.PingTimeout
}

type WebsocketTransport struct {
	PingInterval   time.Duration
	PingTimeout    time.Duration
	ReceiveTimeout time.Duration
	SendTimeout    time.Duration

	BufferSize int

	RequestHeader http.Header
}

func (wst *WebsocketTransport) Connect(url string) (conn transport.Connection, err error) {
	dialer := websocket.Dialer{}

	auth, err := CreateAuthNHeader()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	wst.RequestHeader = make(http.Header)
	wst.RequestHeader.Set("Authorization", auth)
	socket, _, err := dialer.Dial(url, wst.RequestHeader)
	if err != nil {
		return nil, err
	}

	return &WebsocketConnection{socket, wst}, nil
}

func (wst *WebsocketTransport) HandleConnection(
	w http.ResponseWriter, r *http.Request) (conn transport.Connection, err error) {

	if r.Method != "GET" {
		http.Error(w, upgradeFailed+ErrorMethodNotAllowed.Error(), 503)
		return nil, ErrorMethodNotAllowed
	}

	socket, err := websocket.Upgrade(w, r, nil, wst.BufferSize, wst.BufferSize)
	if err != nil {
		http.Error(w, upgradeFailed+err.Error(), 503)
		return nil, errors.New("Http upgrade failed")
	}

	return &WebsocketConnection{socket, wst}, nil
}

/**
Websocket connection do not require any additional processing
*/
func (wst *WebsocketTransport) Serve(w http.ResponseWriter, r *http.Request) {}

/**
Returns websocket connection with default params
*/
func GetDefaultWebsocketTransport() *WebsocketTransport {
	return &WebsocketTransport{
		PingInterval:   60 * time.Second,
		PingTimeout:    30 * time.Second,
		ReceiveTimeout: 60 * time.Second,
		SendTimeout:    60 * time.Second,
		BufferSize:     1024 * 32,
	}
}
