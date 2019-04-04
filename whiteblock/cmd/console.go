package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	util "../util"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var console = &cobra.Command{
	Use:   "console <node>",
	Short: "Logs into the client console",
	Long: `
Console will log into the client console.

Response: stdout of client console`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
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
		log.Fatal(unix.Exec("/usr/bin/ssh", []string{"ssh", "-i", "/home/master-secrets/id.master", "-o", "StrictHostKeyChecking no",
		"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no", "-o", "ConnectTimeout=10", "-y", "-t",
		"root@" + fmt.Sprintf(nodes[nodeNumber].IP), "tmux", "attach", "-t", "whiteblock"}, os.Environ()))
	},
}

func init() {
	RootCmd.AddCommand(console)
}
