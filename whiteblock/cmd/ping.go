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

var pingCmd = &cobra.Command{
	Use:   "ping <sending node> <receiving node>",
	Short: "Ping will send packets to a node.",
	Long: `

Ping will send packets to a node and will output information
Format: <sending node> <receiving node>
Params: sending node, receiving node
	`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			fmt.Println("\nError: Invalid number of arguments given")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "nodes"
		out1 := []byte(wsEmitListen(serverAddr, command1, ""))
		var node Node
		json.Unmarshal(out1, &node)
		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		err = unix.Exec("/usr/bin/ssh", []string{"ssh", "-o", "StrictHostKeyChecking no", "root@" + fmt.Sprintf(node[sendingNodeNumber].IP), "ping", fmt.Sprintf(node[sendingNodeNumber].IP), fmt.Sprintf(node[receivingNodeNumber].IP)}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	pingCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(pingCmd)
}
