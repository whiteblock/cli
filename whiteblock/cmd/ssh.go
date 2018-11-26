package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh <server> <node> <command> ",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		if len(args) != 3 {
			println("Invalid number of arguments given")
			out, err := exec.Command("bash", "-c", "./whiteblock ssh -h").Output()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		command := "exec"
		param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + strings.Join(args[2:], " ") + "\"}"

		wsEmitListen(serverAddr, command, param)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
