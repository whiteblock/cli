package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	util "../util"
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

		util.CheckArguments(args,2,2)
		nodes,err := GetNodes()
		if err != nil{
			util.PrintErrorFatal(err)
		}
		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidArgument(args[0])
			cmd.Help()
			return
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			util.InvalidArgument(args[1])
			cmd.Help()
			return
		}

		if sendingNodeNumber >= len(nodes) {
			util.PrintStringError("Sending node number too high")
			os.Exit(1)
		}

		if receivingNodeNumber >= len(nodes) {
			util.PrintStringError("Receiving node number too high")
			os.Exit(1)
		}

		err = unix.Exec("/usr/bin/ssh", []string{
												"ssh", "-i", "/home/master-secrets/id.master", "-o", "StrictHostKeyChecking no",
												"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no","-o","ConnectTimeout=10", "-y",
												"root@" + fmt.Sprintf(nodes[sendingNodeNumber].IP), "ping", 
												fmt.Sprintf(nodes[receivingNodeNumber].IP)}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
