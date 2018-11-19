package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	command   string
	node      int
	sshserver string
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into an existing container.",
	Long: `SSH will allow the user to go into the contianer where the specified node exists.

	Response: stdout of the command
	`,

	Run: func(cmd *cobra.Command, args []string) {
		//add websocket command later
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		msg := "{\"Servers\":" + fmt.Sprintf("%s", server) + "\",\"node\":" + fmt.Sprintf("%d", node) + ",\"command\":\"" + command + "\"}"

		wsSSH(serverAddr, msg)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "serverAddr", "a", "localhost:5000", "server address with port 5000")
	sshCmd.Flags().StringVarP(&command, "cmd", "c", "ls -l", "Which shell to run in container")
	sshCmd.Flags().StringVarP(&sshserver, "server", "s", "1", "Which server to run in")
	sshCmd.Flags().IntVarP(&node, "node", "n", 0, "Node number to SSH into")
	RootCmd.AddCommand(sshCmd)
}
