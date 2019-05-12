package cmd

import (
	util "../util"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var autoCmd = &cobra.Command{
	Aliases: []string{},
	Use:     "auto <node> <command>",
	Short:   "send queries",
	Long: `Automatically send queries to a node in the background.
`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		node, err := strconv.Atoi(args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		sendPerInterval, err := cmd.Flags().GetInt("send-per-interval")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		jsonRpcCallAndPrint("setup_load", []interface{}{map[string]interface{}{
			"node":            node,
			"name":            fmt.Sprintf("node%d:%s", node, args[1]),
			"interval":        interval,
			"sendPerInterval": sendPerInterval,
			"call":            args[1],
		}})
	},
}

var autoKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill an auto routine",
	Long: `
Kill an auto routine.
`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("state::kill_sub_routines", args)
	},
}

func init() {
	autoCmd.Flags().IntP("interval", "i", 50000, "Send interval in microseconds")
	autoCmd.Flags().IntP("send-per-interval", "b", 1, "Send of requests to send per interval tick (default 1)")
	autoCmd.AddCommand(autoKillCmd)
	RootCmd.AddCommand(autoCmd)
}
