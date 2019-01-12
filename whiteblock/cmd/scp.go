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

var scpCmd = &cobra.Command{
	Use:   "scp <node number> <source> <destination>",
	Short: "Scp will copy a file into the node.",
	Long: `

Scp will allow the user to copy a file and add it to a node.
Format: <node number>, <source>, <destination>
Params: server number, node number, file/dir source, file/dir destination
	`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 3 {
			fmt.Println("\nError: Invalid number of arguments given\n")
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

		err = unix.Exec("/usr/bin/scp", []string{"scp","-i","/home/master-secrets/id.whiteblock", "-r", "-o", "StrictHostKeyChecking no", args[2], "root@" + fmt.Sprintf(node[nodeNumber].IP) + ":" + args[3]}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	scpCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(scpCmd)
}
