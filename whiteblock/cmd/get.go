package cmd

import (
	"fmt"
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

		wsGetServers(serverAddr)
	},
}

var getNodesrCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "List will show all nodes.",
	Long: `
List will output all of the nodes in the current network.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		wsGetNodes(serverAddr)
	},
}

var getTestnetCmd = &cobra.Command{
	Use:   "testnet",
	Short: "Get testnet information",
	Long: `
Testnet will allow the user to get infromation regarding the test network.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		curlGET(fmt.Sprint(serverAddr) + "/testnets/")
	},
}

func init() {
	getCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getCmd.AddCommand(getServerCmd, getTestnetCmd, getNodesrCmd)

	RootCmd.AddCommand(getCmd)
}
