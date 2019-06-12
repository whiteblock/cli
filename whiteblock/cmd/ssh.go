package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"golang.org/x/sys/unix"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	ID        string `json:"id"`
	TestNetID string `json:"testnetId"`
	Server    int    `json:"server"`
	LocalID   int    `json:"localId"`
	IP        string `json:"ip"`
	Label     string `json:"label"`
}

var sshCmd = &cobra.Command{
	Use:   "ssh <node>",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the container where the specified node exists.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, -1)

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
		sshArgs := []string{"ssh", "-i", "/home/master-secrets/id.master", "-o", "StrictHostKeyChecking no",
			"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no", "-o", "ConnectTimeout=10"}
		verbose, err := cmd.Flags().GetBool("verbose")

		if err == nil && verbose {
			sshArgs = append(sshArgs, "-v")
		} else {
			sshArgs = append(sshArgs, "-y")
		}

		sshArgs = append(sshArgs, "root@"+nodes[nodeNumber].IP)

		sshArgs = append(sshArgs, args[1:]...)
		log.WithFields(log.Fields{"command": strings.Join(sshArgs, " ")}).Trace("ssh")
		log.Fatal(unix.Exec("/usr/bin/ssh", sshArgs, os.Environ()))
	},
}

func init() {
	sshCmd.Flags().BoolP("verbose", "v", false, "run ssh in verbose mode")
	RootCmd.AddCommand(sshCmd)
}
