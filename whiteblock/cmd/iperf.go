package cmd

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	bw          string
	testTime    string
	udpEnabled  bool
	dualEnabled bool
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

		CheckArguments(args,2,2)

		nodes,err := GetNodes()
		if err != nil{
			PrintErrorFatal(err)
		}

		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			InvalidArgument(args[0])
			cmd.Help()
			return
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			InvalidArgument(args[1])
			cmd.Help()
			return
		}

		wg.Add(2)
		// command to run iperf as a server
		go func() {
			defer wg.Done()

			iPerfcmd := "iperf3 -s "
			if udpEnabled {
				iPerfcmd = iPerfcmd + "-u "
			}

			iPerfcmd = iPerfcmd + fmt.Sprintf(nodes[sendingNodeNumber].IP) + " -1"

			client, err := NewSshClient(fmt.Sprintf(nodes[sendingNodeNumber].IP))
			if err != nil {
				panic(err)
			}
			defer client.Close()

			client.Run("pkill -9 iperf3")

			result, err := client.Run(iPerfcmd)
			if err != nil {
				fmt.Println(result)
				panic(err)
			}
			fmt.Println(result)

		}()

		go func() {
			// command to run iperf as a client
			time.Sleep(5 * time.Second)
			defer wg.Done()

			iPerfcmd := "iperf3 -c "
			if udpEnabled {
				iPerfcmd = iPerfcmd + " -u "
			}
			if bw != "" && udpEnabled {
				_, err := strconv.Atoi(bw)
				if err != nil {
					fmt.Println("Invalid format given for bandwidth flag.")
					return
				}
				iPerfcmd = iPerfcmd + " -b " + bw
			} else if bw != "" && !udpEnabled {
				fmt.Println("udp needs to be enabled to set bandwidth.")
			}
			if dualEnabled {
				iPerfcmd = iPerfcmd + " -d "
			}

			iPerfcmd = iPerfcmd + fmt.Sprintf(nodes[sendingNodeNumber].IP)

			client, err := NewSshClient(fmt.Sprintf(nodes[receivingNodeNumber].IP))
			if err != nil {
				panic(err)
			}
			defer client.Close()

			client.Run("pkill -9 iperf3")

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
	iPerfCmd.Flags().StringVarP(&bw, "bandwidth", "b", "", "set target bandwidth in bits/sec (default 1 Mbit/sec); requires udp enabled")
	iPerfCmd.Flags().BoolVarP(&dualEnabled, "dualtest", "d", false, "enable bidirectional test simultaneously")
	iPerfCmd.Flags().StringVarP(&testTime, "time", "t", "", "how long to run test for")
	iPerfCmd.Flags().BoolVarP(&udpEnabled, "udp", "u", false, "enable udp")

	RootCmd.AddCommand(iPerfCmd)
}
