package cmd

import (
	util "../util"
	"github.com/spf13/cobra"
)

var restartNodeCmd = &cobra.Command{
	Use:   "restart [node number]",
	Short: "Attempt to restart a node",
	Long: `
Kill a node by sending SIGINT and then re-run the original command used to run it`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}

		jsonRpcCallAndPrint("restart_node", []interface{}{testnetId, args[0]})
	},
}

func init() {
	RootCmd.AddCommand(restartNodeCmd)
}
