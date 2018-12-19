package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

func BashExec(_cmd string) (string, error) {
	cmd := exec.Command("bash", "-c", _cmd)
	var resultsRaw bytes.Buffer
	cmd.Stdout = &resultsRaw
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	return resultsRaw.String(), nil
}

var iPerfCmd = &cobra.Command{
	Use:   "iperf <sending node> <receiving node>",
	Short: "iperf will show network conditions.",
	Long: `

Iperf will show the user network conditions and other data.

Format: <sending node> <receiving node>
Params: sending node, receiving node
	`,

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		if len(args) != 2 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "nodes"
		out1 := []byte(wsEmitListen(serverAddr, command1, ""))
		var node Node
		json.Unmarshal(out1, &node)
		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		command2 := "exec"
		param := "{\"server\":5,\"node\":" + args[0] + ",\"command\":\"service ssh start\"}"
		wsEmitListen(serverAddr, command2, param)

		go func() {
			wg.Add(1)
			defer wg.Done()
			iPerfcmd1 := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "root@"+fmt.Sprintf(node[sendingNodeNumber].IP), "iperf3", fmt.Sprintf(node[sendingNodeNumber].IP), fmt.Sprintf(node[receivingNodeNumber].IP))
			err = iPerfcmd1.Run()
			if err != nil {
				panic(err)
			}
		}()

		go func() {
			wg.Add(1)
			defer wg.Done()
			iPerfcmd2 := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "root@"+fmt.Sprintf(node[receivingNodeNumber].IP), "iperf3", fmt.Sprintf(node[receivingNodeNumber].IP), fmt.Sprintf(node[sendingNodeNumber].IP))
			err := iPerfcmd2.Run()
			if err != nil {
				panic(err)
			}

		}()

		wg.Wait()
	},
}

func init() {
	iPerfCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(iPerfCmd)
}
