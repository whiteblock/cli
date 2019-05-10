package cmd

import (
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Hidden: true,
	Use:    "exec",
	Short:  "Execute a function call",
	Long:   "\nMainly for internal and debug purposes.\n",
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint(args[0], args[1:])
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
}
