package cmd

import (
	"github.com/spf13/cobra"
)

var (
	servers string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Get server information.",
	Long: `Server will allow the user to get server information.
	
	info				Get the information from all currently registered servers;
	`,
	// info id [Server ID]		Get server information by id;

	Run: func(cmd *cobra.Command, args []string) {
		// curlGET(fmt.Sprint(serverAddr) + "/servers/" + fmt.Sprint(servers))
		// msg := "get_servers"
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		wsGetServers(serverAddr)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&servers, "ID", "i", "", "Server ID")
	serverCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(serverCmd)
}
