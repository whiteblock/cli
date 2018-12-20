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
		param := "{\"server\":1,\"node\":" + args[0] + ",\"command\":\"service ssh start\"}"
		wsEmitListen(serverAddr, command2, param)

		wg.Add(2)
		go func() {
			defer wg.Done()
			// iPerfcmd1 := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "root@"+fmt.Sprintf(node[sendingNodeNumber].IP), "iperf3", "-s", fmt.Sprintf(node[sendingNodeNumber].IP))
			// err = iPerfcmd1.Run()
			// if err != nil {
			// 	fmt.Println(err)
			// }

			// writer := bufio.NewWriterSize(os.Stdout, 20)
			// go func() {
			// 	for {
			// 		if writer.Available() == 10 {
			// 			writer.Flush()
			// 		}
			// 	}
			// }()
			// iPerfcmd1.Stdout = writer
			// err = iPerfcmd1.Start()
			// if err != nil {
			// 	panic(err)
			// }
			// err = iPerfcmd1.Wait()
			// if err != nil {
			// 	panic(err)
			// }
			// writer.Flush()

			iPerfcmd := "iperf3 -s " + fmt.Sprintf(node[sendingNodeNumber].IP) + " -1"

			client, err := NewSshClient(fmt.Sprintf(node[sendingNodeNumber].IP))
			if err != nil {
				panic(err)
			}
			result, err := client.Run(iPerfcmd)
			if err != nil {
				fmt.Println(result)
				panic(err)
			}
			fmt.Println(result)

		}()

		go func() {
			time.Sleep(3 * time.Second)
			defer wg.Done()
			iPerfcmd := "iperf3 -c " + fmt.Sprintf(node[sendingNodeNumber].IP)

			client, err := NewSshClient(fmt.Sprintf(node[receivingNodeNumber].IP))
			if err != nil {
				panic(err)
			}
			result, err := client.Run(iPerfcmd)
			if err != nil {
				fmt.Println(result)
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
