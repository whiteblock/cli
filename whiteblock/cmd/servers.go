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
	
	serverInfo				Get the information from all currently registered servers;
	serverInfo --id [Server ID]		Get server information by id;
	`,

	Run: func(cmd *cobra.Command, args []string) {
		curlGET("http://localhost:8000/servers/" + fmt.Sprint(servers))
	},
}

func init() {
	serverCmd.Flags().StringVarP(&servers, "ID", "i", "", "Server ID")

	RootCmd.AddCommand(serverCmd)
}
