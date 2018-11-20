package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh <node> <server> <command> ",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command
	`,

	Run: func(cmd *cobra.Command, args []string) {
		//add websocket command later
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		// msg := "{\"server\":" + fmt.Sprintf("%s", sshserver) + ",\"node\":" + fmt.Sprintf("%d", node) + ",\"command\":\"" + command + "\"}"

		if len(args) < 3 {
			println("Invalid number of arguments given")
			out, err := exec.Command("bash", "-c", "./whiteblock ssh -h").Output()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", out)
		}

		msg := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + args[2] + "\"}"

		wsSSH(serverAddr, msg)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
