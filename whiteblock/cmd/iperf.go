package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var iPerfCmd = &cobra.Command{
	Use:   "iperf <sending node> <receiving node>",
	Short: "iperf will show network conditions.",
	Long: `

Iperf will show the user network conditions and other data. This command will establish the sending node as a server and the receiving node as a client node. They will send packets and at the end of the test, the output will give bandwidth, transfer size, and other relevant

Format: <sending node> <receiving node>
Params: sending node, receiving node
	`,

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		if len(args) != 2 {
			fmt.Println("\nError: Invalid number of arguments given")
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
		
		wg.Add(2)
		go func() {
			defer wg.Done()

			iPerfcmd := "iperf3 -s " + fmt.Sprintf(node[sendingNodeNumber].IP) + " -1"

			client, err := NewSshClient(fmt.Sprintf(node[sendingNodeNumber].IP))
			if err != nil {
				panic(err)
			}
			defer client.Close()
			result, err := client.Run(iPerfcmd)
			if err != nil {
				fmt.Println(result)
				panic(err)
			}
			fmt.Println(result)

		}()

		go func() {
			time.Sleep(5 * time.Second)
			defer wg.Done()
			iPerfcmd := "iperf3 -c " + fmt.Sprintf(node[sendingNodeNumber].IP)

			client, err := NewSshClient(fmt.Sprintf(node[receivingNodeNumber].IP))
			if err != nil {
				panic(err)
			}
			defer client.Close()
			result, err := client.Run(iPerfcmd)
			if err != nil {
				fmt.Println(result)
				panic(err)
			}
			fmt.Println(result)
		}()

		wg.Wait()
	},
}

func init() {
	iPerfCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(iPerfCmd)
}
