package cmd

import (
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
		CheckArguments(args,3,3)

		nodes,err := GetNodes()
		if err != nil{
			PrintErrorFatal(err)
		}
		nodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		err = unix.Exec("/usr/bin/scp", []string{"scp","-i","/home/master-secrets/id.master", "-r", "-o","UserKnownHostsFile=/dev/null",
						"-o", "StrictHostKeyChecking no", args[1], "root@" + fmt.Sprintf(nodes[nodeNumber].IP) + ":" + args[2]}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	scpCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(scpCmd)
}
