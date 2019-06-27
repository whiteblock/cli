package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
)

var syncCmd = &cobra.Command{
	Aliases: []string{"pull"},
	Use:     "sync",
	Short:   "Sync up with your current state",
	Long: `
	Sync up with your current state.
`,

	Run: func(cmd *cobra.Command, args []string) {

		res, err := util.JsonRpcCall("get_last_build", []interface{}{})
		if err != nil {
			util.PrintErrorFatal(err)
		}
		err = util.WriteStore(".previous_build_id", []byte(res.(map[string]interface{})["id"].(string)))
		if err != nil {
			util.PrintErrorFatal(err)
		}
		cmd.Println("synced up with the latest build")
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
