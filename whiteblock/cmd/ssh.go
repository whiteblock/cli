package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var (
	dir     string
	cwd     []string
	lastcmd string
)

type Node []struct {
	ID        int    `json:"id"`
	TestNetID int    `json:"testNetId"`
	Server    int    `json:"server"`
	LocalID   int    `json:"localId"`
	IP        string `json:"ip"`
	Label     string `json:"label"`
}

var sshCmd = &cobra.Command{
	Use:   "ssh <server> <node>",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command
	`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "nodes"
		out1 := []byte(wsEmitListen(serverAddr, command1, ""))
		var node Node
		json.Unmarshal(out1, &node)
		nodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		println(node[nodeNumber].IP)

		command2 := "exec"
		param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"service ssh start\"}"

		println(wsEmitListen(serverAddr, command2, param))

		err = unix.Exec("/usr/bin/ssh", []string{"ssh", "root@" + fmt.Sprintf(node[nodeNumber].IP)}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
