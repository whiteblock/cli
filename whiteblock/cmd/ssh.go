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

type Node struct {
	ID        string `json:"id"`
	TestNetID string `json:"testNetId"`
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
		if len(args) < 1 {
			fmt.Println("\nError: Invalid number of arguments given")
			cmd.Help()
			return
		}

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
			"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no", "-y",
			"root@" + fmt.Sprintf(nodes[nodeNumber].IP)}

		sshArgs = append(sshArgs, args[1:]...)
		//fmt.Println(strings.Join(sshArgs," "))
		log.Fatal(unix.Exec("/usr/bin/ssh", sshArgs, os.Environ()))
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
