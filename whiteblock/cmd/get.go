package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <command>",
	Short: "Get server and network information.",
	Long: `
Get will allow the user to get server and network information and statstics.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock get -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var getServerCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"servers"},
	Short:   "Get server information.",
	Long: `
Server will allow the user to get server information.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_servers"

		wsEmitListen(serverAddr, command, "")
	},
}

var getNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Nodes will show all nodes.",
	Long: `
Nodes will output all of the nodes in the current network.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_nodes"
		wsEmitListen(serverAddr, command, "")
	},
}

var getStatsCmd = &cobra.Command{
	Use:   "stats <command>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock get stats -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var statsByTimeCmd = &cobra.Command{
	Use:   "time <start time> <end time>",
	Short: "Get stastics by time",
	Long: `
Stats will allow the user to get statistics by specifying a start time and stop time (unix time stamp).

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":" + args[0] + ",\"endTime\":" + args[1] + ",\"startBlock\":0,\"endBlock\":0}"
		wsEmitListen(serverAddr, command, param)
	},
}

var statsByBlockCmd = &cobra.Command{
	Use:   "block <start block> <end block>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Params: Block numbers
Format: <start block number> <end block number>

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":0,\"endTime\":0,\"startBlock\":" + args[0] + ",\"endBlock\":" + args[1] + "}"
		wsEmitListen(serverAddr, command, param)
	},
}

func init() {
	getCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd)
	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd)

	RootCmd.AddCommand(getCmd)
}
