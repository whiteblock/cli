package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
)

var signalCmd = &cobra.Command{
	Aliases: []string{"raise"},
	Use:     "signal <node> [sig=SIGTERM]",
	Short:   "Raise a signal to a node's main process",
	Long:    `Sends a signal to the node's main process, see signal(7) for more details about signal`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		signal := "SIGTERM"
		if len(args) > 1 {
			signal = args[1]
		}

		util.JsonRpcCallAndPrint("signal_node", []interface{}{testnetID, args[0], signal})
	},
}

func init() {
	RootCmd.AddCommand(signalCmd)
}
