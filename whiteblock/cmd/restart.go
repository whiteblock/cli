package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
)

var restartNodeCmd = &cobra.Command{
	Use:   "restart [node number]",
	Short: "Attempt to restart a node",
	Long: `
Kill a node by sending SIGINT and then re-run the original command used to run it`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		util.JsonRpcCallAndPrint("restart_node", []interface{}{build.GetPreviousBuildID(), args[0]})
	},
}

func init() {
	RootCmd.AddCommand(restartNodeCmd)
}
