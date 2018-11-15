package cmd

import (
	"github.com/spf13/cobra"
)

var (
	command string
	node    int
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into an existing container.",
	Long: `SSH will allow the user to go into the contianer where the specified node exists.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		//add websocket command later
	},
}

func init() {
	sshCmd.Flags().StringVarP(&command, "cmd", "c", "bash", "Which shell to run in container")
	sshCmd.Flags().IntVarP(&node, "node", "n", 0, "Node number to SSH into")

	RootCmd.AddCommand(sshCmd)
}
