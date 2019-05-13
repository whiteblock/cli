package cmd

import (
	util "../util"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var autoCmd = &cobra.Command{
	Aliases: []string{},
	Use:     "auto <node> <command> [params]",
	Short:   "send queries",
	Long: `Automatically send queries to a node in the background. <command> is the name of the json rpc call to be made
`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, -1)
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
		params := []interface{}{}
		if len(args) > 2 {
			for _, arg := range args[2:] {
				var param interface{}
				err = json.Unmarshal([]byte(arg), &param)
				if err != nil {
					param = arg //if it is not json, then it is a string
				}
				params = append(params, param)
			}
		}
		jsonRpcCallAndPrint("setup_load", []interface{}{map[string]interface{}{
			"node":            node,
			"name":            fmt.Sprintf("node%d:%s", node, args[1]),
			"interval":        interval,
			"sendPerInterval": sendPerInterval,
			"call":            args[1],
			"arguments":       params,
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
