package cmd

import (
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

		CheckArguments(args,2,2)
		nodes,err := GetNodes()
		if err != nil{
			PrintErrorFatal(err)
		}
		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			InvalidArgument(args[0])
			cmd.Help()
			return
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			InvalidArgument(args[1])
			cmd.Help()
			return
		}
		err = unix.Exec("/usr/bin/ssh", []string{"ssh","-i","/home/master-secrets/id.master",
												"-o","UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking no", 
												"root@" + fmt.Sprintf(nodes[sendingNodeNumber].IP), "ping", 
												 fmt.Sprintf(nodes[receivingNodeNumber].IP)}, os.Environ())
		log.Fatal(err)
	},
}

func init() {
	pingCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(pingCmd)
}
