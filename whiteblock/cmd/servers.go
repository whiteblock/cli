package cmd

import (
	"fmt"

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
		msg := "{\"Servers\":" + fmt.Sprintf("v", server) + ",\"Blockchain\":" + blockchain + ",\"Nodes\":" + fmt.Sprintf("%d", nodes) + ",\"Image\":" + image + "}"

		wsEmit(serverAddr, msg)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&servers, "ID", "i", "", "Server ID")
	serverCmd.Flags().StringVarP(&serverAddr, "serverAddr", "a", "http://localhost:8000", "server address with port 8000")

	RootCmd.AddCommand(serverCmd)
}
