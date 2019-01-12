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

type Node []struct {
	ID        int    `json:"id"`
	TestNetID int    `json:"testNetId"`
	Server    int    `json:"server"`
	LocalID   int    `json:"localId"`
	IP        string `json:"ip"`
	Label     string `json:"label"`
}

var sshCmd = &cobra.Command{
	Use:   "ssh <node>",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.
	`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("\nError: Invalid number of arguments given")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "nodes"
		out1 := []byte(wsEmitListen(serverAddr, command1, ""))
		var node Node
		json.Unmarshal(out1, &node)
		nodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}

		log.Fatal(unix.Exec("/usr/bin/ssh", []string{"ssh","-i","/home/master-secrets/id.customer", "-o", "StrictHostKeyChecking no", 
			"-o","UserKnownHostsFile=/dev/null","root@" + fmt.Sprintf(node[nodeNumber].IP)}, os.Environ()))
		fmt.Println(nodeNumber)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
