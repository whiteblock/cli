package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
)

var execCmd = &cobra.Command{
	Hidden: true,
	Use:    "exec",
	Short:  "Execute a function call",
	Long:   "\nMainly for internal and debug purposes.\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint(args[0], util.ArgsToJSON(args[1:]))
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
}
