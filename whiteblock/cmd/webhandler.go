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


func buildListener(){
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

	c.Emit("build_status", "")
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

	// build
	if cmd == "build" {
		err = c.On("build", func(h *gosocketio.Channel, args string) {
			log.Println("build: ", args)
			if args != "Build in Progress" {
				os.Exit(1)
			}
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = c.On("build_status", func(h *gosocketio.Channel, args string) {
			var status BuildStatus
			json.Unmarshal([]byte(args), &status)
			if status.Progress == 0.0 {
				fmt.Printf("Sending build context to Whiteblock\r")
			} else {
                fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\r",status.Stage, status.Progress)
			}

			if status.Error != nil {
				what := status.Error["what"]
				fmt.Println("\n" + what)
				mutex.Unlock()
			} else if status.Progress == 100.0 {
                fmt.Println("\a")
				mutex.Unlock()
			}
		})
	}

	// eos commands
	if strings.HasPrefix(cmd, "eos::") {
		err = c.On(cmd, func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = prettyp(args)
			}
			mutex.Unlock()
		})
	}

	// gethcmd
	if strings.HasPrefix(cmd, "eth::") {
		if cmd == "eth::start_mining" {
			fmt.Println("Started mining. Please wait for the DAG to be generated. Number of miners started: ")
		} else if cmd == "eth::stop_mining" {
			fmt.Println("Stopped mining. Number of miners stopped: ")
		} else if cmd == "eth::start_transactions" {
			fmt.Println("Please make sure that mining has started for the transactions to be included in the blocks.")
		}

		err = c.On(cmd, func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = prettyp(args)
				mutex.Unlock()
			}
		})
	}

	// get_defaults
	if cmd == "get_defaults" {
		err = c.On("get_defaults", func(h *gosocketio.Channel, args string) {
			out = prettyp(args)
			mutex.Unlock()
		})
	}

	// get_params
	if cmd == "get_params" {
		err = c.On("get_params", func(h *gosocketio.Channel, args string) {
			out = args
			mutex.Unlock()
		})
	}

	// get servers
	if cmd == "get_servers" {
		err = c.On("get_servers", func(h *gosocketio.Channel, args string) {
			out = (args)
			mutex.Unlock()
		})
	}

	// get nodes
	if cmd == "get_nodes" {
		err = c.On("get_nodes", func(h *gosocketio.Channel, args string) {
			out = prettyp(args)
			mutex.Unlock()
		})
	}

	// get stats
	if cmd == "stats" {
		err = c.On("stats", func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = prettyp(args)
			} else {
				fmt.Println(err.Error())
			}
			mutex.Unlock()
		})
	}

	// get log
	if cmd == "log" {
		err = c.On("log", func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = args
			} else {
				fmt.Println("Not supported!")
			}
			mutex.Unlock()
		})
	}

	// nodes
	if cmd == "nodes" {
		err = c.On("nodes", func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = args
			} else {
				fmt.Println("Command failed")
			}
			mutex.Unlock()
		})
	}

	// get all_stats
	if cmd == "all_stats" {
		err = c.On("all_stats", func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = prettyp(args)
			} else {
				fmt.Println(err.Error())
			}
			mutex.Unlock()
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
			out = args
			mutex.Unlock()
		})
	}

	// state
	if strings.HasPrefix(cmd, "state::") {
		err = c.On(cmd, func(h *gosocketio.Channel, args string) {
			if len(args) > 0 {
				out = args
				mutex.Unlock()
			}
		})
	}

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
