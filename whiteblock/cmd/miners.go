package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var minerCmd = &cobra.Command{
	// Hidden: true,
	Use:   "miner <command>",
	Short: "Run miner commands.",
	Long: `
Send commands pertaining to mining. This will be blockchain specific and will only be supported depending on which blockchain had been built.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.")
		cmd.Help()
		return
	},
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := ""
		param := ""
		switch blockchain {
		case "ethereum":
			command = "eth::start_mining"
			param = strings.Join(args[:], " ")
			out := fmt.Sprintf("%s", wsEmitListen(serverAddr, command, param))

			if out == "ERROR" {
				fmt.Println("There was an error building the DAG.")
				return
			}
			DagReady := false
			for !DagReady {
				fmt.Printf("\rDAG is being generated...")
				blocknum, _ := strconv.Atoi(wsEmitListen(serverAddr, "eth::get_block_number", ""))
				if blocknum > 4 {
					DagReady = true
				}
			}
			fmt.Println("\rDAG has been successfully generated.")
		case "eos":
			fmt.Println("This function is not supported for the eos client.")
			return
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := ""
		param := ""
		switch blockchain {
		case "ethereum":
			command = "eth::stop_mining"
			param = strings.Join(args[:], " ")
		case "eos":
			fmt.Println("This function is not supported for the eos client.")
			return
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

func init() {
	minerCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	minerCmd.AddCommand(minerStartCmd, minerStopCmd)
	RootCmd.AddCommand(minerCmd)
}
