package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"golang.org/x/sys/unix"
	"log"
	"os"
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
		nodes := GetNodes()
		nodeNumber := util.CheckAndConvertInt(args[0], "node")
		util.CheckIntegerBounds(cmd, "node number", nodeNumber, 0, len(nodes)-1)

		log.Fatal(unix.Exec("/usr/bin/scp", []string{"scp", "-i", conf.SSHPrivateKey,
			"-r", "-o", "UserKnownHostsFile=/dev/null",
			"-o", "StrictHostKeyChecking no", args[1],
			"root@" + fmt.Sprintf(nodes[nodeNumber].IP) + ":" + args[2]}, os.Environ()))
	},
}

func init() {
	RootCmd.AddCommand(scpCmd)
}
