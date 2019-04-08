package cmd

import (
	util "../util"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Aliases: []string{"die", "stop", "teardown"},
	Use:     "done",
	Short:   "Tears down the testnet",
	Long: `
	Tears down the nodes, and frees up any resources which they are using.
`,

	Run: func(cmd *cobra.Command, args []string) {
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		jsonRpcCallAndPrint("delete_testnet", []interface{}{testnetId})
	},
}

func init() {
	RootCmd.AddCommand(doneCmd)
}
