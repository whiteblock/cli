package cmd

import (
	util "../util"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Aliases: []string{},
	Use:     "kill <node>",
	Short:   "Raise SIGINT to a node's main process and wait for it to die",
	Long: `Sends SIGINT to the node's main process, and continue to query the state of that 
	process until it dies. 
`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		jsonRpcCallAndPrint("kill_node", []interface{}{testnetId, args[0]})
	},
}

func init() {
	RootCmd.AddCommand(killCmd)
}
