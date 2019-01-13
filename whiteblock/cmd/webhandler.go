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
    Error       map[string]string   `json:"error"`
    Progress    float64             `json:"progress"`
    Stage       string              `json:"stage"`
}

func wsEmitListen(wsaddr, cmd, param string) string {
    var mutex = &sync.Mutex{}
    mutex.Lock()
    c, err := gosocketio.Dial(
        wsaddr,
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
            if len(status.Stage) == 0 {
                fmt.Printf("Sending build context to Whiteblock\r")
            }else{
                fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\a\r",status.Stage, status.Progress)
            }
            
            if status.Error != nil {
                what := status.Error["what"]
                fmt.Println("\n" + what)
                mutex.Unlock()
            } else if status.Progress == 100.0 {
                fmt.Println("")
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
        err = c.On(cmd, func(h *gosocketio.Channel, args string) {
            if len(args) > 0 {
                if strings.Contains(args, "{") {
                    fmt.Println(prettyp(args))
                } else {
                    fmt.Println(args)
                }
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
                fmt.Println(err.Error())
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
                fmt.Println(err.Error())
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
