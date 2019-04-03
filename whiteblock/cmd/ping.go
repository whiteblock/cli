package cmd

import (
	util "../util"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"strconv"
)

var pingCmd = &cobra.Command{
	Use:   "ping <sending node> <receiving node>",
	Short: "Ping will send packets to a node.",
	Long: `
Ping will send packets to a node and will output information

Params: sending node, receiving node
	`,

	Run: func(cmd *cobra.Command, args []string) {

		util.CheckArguments(cmd, args, 2, 2)
		nodes, err := GetNodes()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidInteger("sending node number", args[0], true)
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			util.InvalidInteger("receiving node number", args[1], true)
		}
		util.CheckIntegerBounds(cmd, "sending node number", sendingNodeNumber, 0, len(nodes)-1)
		util.CheckIntegerBounds(cmd, "receiving node number", receivingNodeNumber, 0, len(nodes)-1)

		err = unix.Exec("/usr/bin/ssh", []string{
			"ssh", "-i", "/home/master-secrets/id.master", "-o", "StrictHostKeyChecking no",
			"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no", "-o", "ConnectTimeout=10", "-y",
			"root@" + fmt.Sprintf(nodes[sendingNodeNumber].IP), "ping",
			fmt.Sprintf(nodes[receivingNodeNumber].IP)}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
