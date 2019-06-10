package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"strconv"
)

var scpCmd = &cobra.Command{
	Use:   "scp <node> <source> <destination>",
	Short: "Scp will copy a file into the node.",
	Long: `

Scp will allow the user to copy a file and add it to a node.
Format: <node>, <source>, <destination>
Params: node number, file/dir source, file/dir destination
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 3, 3)

		nodes, err := GetNodes()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		nodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if nodeNumber >= len(nodes) {
			util.PrintStringError("Node number too high")
			os.Exit(1)
		}
		err = unix.Exec("/usr/bin/scp", []string{"scp", "-i", "/home/master-secrets/id.master", "-r", "-o", "UserKnownHostsFile=/dev/null",
			"-o", "StrictHostKeyChecking no", args[1], "root@" + fmt.Sprintf(nodes[nodeNumber].IP) + ":" + args[2]}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	RootCmd.AddCommand(scpCmd)
}
