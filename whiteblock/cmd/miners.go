package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

var minerCmd = &cobra.Command{
	// Hidden: true,
	Use:   "miner <command>",
	Short: "Run miner commands.",
	Long: `
Send commands pertaining to mining. This will be blockchain specific and will only be supported depending on which blockchain had been built.
	`,
	Run: PartialCommand,
}

var minerStartCmd = &cobra.Command{
	Use:   "start [node 1 number] [node 2 number]...",
	Short: "Start Mining",
	Long: `
Send the start mining signal to nodes, may take a while to take effect due to DAG generation. If no arguments are given, all nodes will begin mining.

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to start mining or None for all nodes

Response: The number of nodes which successfully received the signal to start mining`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(blockchain) == 0 {
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		switch blockchain {
			case "ethereum":
				res,err := jsonRpcCall("eth::start_mining",args)
				if err != nil {
					PrintStringError(err.Error())
					PrintStringError("There was an error building the DAG.")
					os.Exit(1)
				}
				DagReady := false
				for !DagReady {
					fmt.Printf("\rDAG is being generated...")
					res,err = jsonRpcCall("eth::get_block_number",[]string{})
					if err != nil {
						PrintErrorFatal(err)
					}
					blocknum := int(res.(float64))
					if blocknum > 4 {
						DagReady = true
					}
				}
				fmt.Println("\rDAG has been successfully generated.")
			default:
				ClientNotSupported(blockchain)
		}
	},
}

var minerStopCmd = &cobra.Command{
	Use:   "stop [node 1 number] [node 2 number]...",
	Short: "Stop mining",
	Long: `
Send the stop mining signal to nodes

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to stop mining or None for all nodes


Response: The number of nodes which successfully received the signal to stop mining`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(blockchain) == 0 {
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		command := ""
		switch blockchain {
			case "ethereum":
				command = "eth::stop_mining"
			default:
				ClientNotSupported(blockchain)
		}
		jsonRpcCallAndPrint(command,args)
	},
}

func init() {
	minerCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	minerCmd.AddCommand(minerStartCmd, minerStopCmd)
	RootCmd.AddCommand(minerCmd)
}
