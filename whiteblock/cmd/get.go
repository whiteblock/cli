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
Get will allow the user to get server and network information.
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

var getTestnetCmd = &cobra.Command{
	Use:   "testnet",
	Short: "Get testnet information",
	Long: `
Testnet will allow the user to get infromation regarding the test network.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		// add get_testnet command
	},
}

func init() {
	getCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getCmd.AddCommand(getServerCmd, getTestnetCmd, getNodesCmd)

	RootCmd.AddCommand(getCmd)
}
