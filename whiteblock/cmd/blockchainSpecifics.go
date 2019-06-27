package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
)

var eosCmd = &cobra.Command{
	Use:   "eos <command>",
	Short: "Run eos commands",
	Long:  "\nEos will allow the user to get information and run EOS commands.\n",
	Run:   util.PartialCommand,
}

var eosGetInfoCmd = &cobra.Command{
	Use:   "get_info [node]",
	Short: "Get EOS info",
	Long: `
Roughly equivalent to calling cleos get info

Params: The node to get info from

Response: eos blockchain state info`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		util.JsonRpcCallAndPrint("eos::get_info", args)
	},
}

func init() {
	eosCmd.AddCommand(eosGetInfoCmd)
	RootCmd.AddCommand(eosCmd)
}
