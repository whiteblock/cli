package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"golang.org/x/sys/unix"
	"log"
	"os"
)

var console = &cobra.Command{
	Use:   "console <node>",
	Short: "Logs into the client console",
	Long: `
Console will log into the client console.

Response: stdout of client console`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		nodes := GetNodes()
		nodeNumber := util.CheckAndConvertInt(args[0], "node")
		util.CheckIntegerBounds(cmd, "node number", nodeNumber, 0, len(nodes)-1)

		log.Fatal(unix.Exec(conf.SSHBinary, []string{"ssh", "-i", conf.SSHPrivateKey, "-o", "StrictHostKeyChecking no",
			"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no", "-o", "ConnectTimeout=10", "-y", "-t",
			"root@" + fmt.Sprintf(nodes[nodeNumber].IP), "tmux", "attach", "-t", "whiteblock"}, os.Environ()))
	},
}

func init() {
	RootCmd.AddCommand(console)
}
