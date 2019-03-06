package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"net/http"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	util "../util"
)

type BuildStatus struct {
	Error    map[string]string `json:"error"`
	Progress float64           `json:"progress"`
	Stage    string            `json:"stage"`
}

func jsonRpcCallAndPrint(method string,params interface{}) {
	reply,err := jsonRpcCall(method,params)
	switch reply.(type) {
		case string:
			fmt.Printf("\033[97m%s\033[0m\n",reply.(string))
			return
	}
	if err != nil {
		util.PrintErrorFatal(err)
	}
	fmt.Println(prettypi(reply))
}

func jsonRpcCall(method string,params interface{}) (interface{},error) {
	//log.Println("URL IS "+url)
	jrpc,err := json2.EncodeClientRequest(method,params)
	if err != nil{
		log.Println(err)
		return nil,err
	}
	body := strings.NewReader(string(jrpc))
	req, err := http.NewRequest("POST",fmt.Sprintf("http://%s/rpc",serverAddr),body)
	if err != nil {
		log.Println(err)
		return nil,err
	}
	req.Header.Set("Content-Type", "application/json")

	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil,err
	}
	defer resp.Body.Close()
	var out interface{}
	err = json2.DecodeClientResponse(resp.Body,&out)
	if err != nil {
		//log.Println(err)
		return nil,err
	}
	return out, nil
}


func buildListener(testnetId string){
	var mutex = &sync.Mutex{}
	mutex.Lock()
	c, err := gosocketio.Dial(
		"ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket",
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	defer c.Close()

	c.On("error",func(h *gosocketio.Channel, args string){
		util.PrintStringError(args)
		os.Exit(1)
	})

	
	err = c.On("build_status", func(h *gosocketio.Channel, args string) {
		var status BuildStatus
		json.Unmarshal([]byte(args), &status)
		if status.Progress == 0.0 {
			fmt.Printf("Sending build context to Whiteblock\r")
		}else if status.Error != nil {
			what := status.Error["what"]
			util.PrintStringError(what)
			mutex.Unlock()
		} else if status.Progress == 100.0 {
            fmt.Println("\a")
			mutex.Unlock()
		} else {
            fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\r",status.Stage, status.Progress)
		}
	})

	c.Emit("build_status", testnetId)
	mutex.Lock()
}

func wsEmitListen(wsaddr, cmd, param string) string {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	c, err := gosocketio.Dial(
		"ws://" + wsaddr + "/socket.io/?EIO=3&transport=websocket",
		transport.GetDefaultWebsocketTransport(),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer c.Close()

	out := ""

	// sys commands
	if strings.HasPrefix(cmd, "sys::") {
		err = c.On(cmd, func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = args
				mutex.Unlock()
			}
		})
	}

	c.Emit(cmd, param)
	mutex.Lock()
	return out
}
