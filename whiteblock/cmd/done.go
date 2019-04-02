package cmd

import (
	util "../util"
	"github.com/spf13/cobra"
)


var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Tear down the testnet",
	Long: `
	Tears down the in use resource of the built testnet
`,

	Run: func(cmd *cobra.Command, args []string) {
		testnetId,err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		jsonRpcCallAndPrint("delete_testnet", []interface{}{testnetId})
	},
}


func init() {	
	RootCmd.AddCommand(doneCmd)
}
