package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"os"
)

var jsonrpcCall = &cobra.Command{
	Use:   "jsonrpc <node> <command> [args..]",
	Short: "send a json rpc call",
	Long:  "\nSend a json rpc call to a node.\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, util.NoMaxArgs)
		nodes, err := GetNodes()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		nodeNumber := util.CheckAndConvertInt(args[0], "node number")

		if nodeNumber >= len(nodes) {
			util.PrintStringError("Node number too high")
			os.Exit(1)
		}
		util.JsonRpcCallAndPrint("jsonrpc_call", args)
	},
}

func init() {
	RootCmd.AddCommand(jsonrpcCall)
}
