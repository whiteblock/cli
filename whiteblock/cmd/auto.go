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
	Long: `Automatically send json_rpc queries to a node in the background. <command> is the name of the json rpc call to be made. 
	You can use +account,+tx_hash,+number,+hex,+block_hash,+block_number as magic string parameters to be filled in with randomized appropriate values.
	+tx_hash random tx hash; only works after you call wb tx start stream
	+account random account
	+number random base 10 number
	+hex random hex number
	+block_hash random block hash
	+block_number random block number
	Examples:
	wb auto 0 eth_sendTransaction -i 1000000 '{"from":"+account","to":"+account","gas":"0x76c0","gasPrice":"0x9184e72a000","value":"+hex","data":"0x00"}'
	wb auto 0 eth_getBalance -i 100000 +account latest
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
		errorChecking, err := cmd.Flags().GetBool("full-error-checking")
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
			"errorCheck":      errorChecking,
		}})
	},
}

var autoKillCmd = &cobra.Command{
	Use:     "kill",
	Aliases: []string{"stop"},
	Short:   "stop an auto routine",
	Long: `
Kill an auto routine.
`,
	Run: func(cmd *cobra.Command, args []string) {
		forced, err := cmd.Flags().GetBool("force")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if forced {
			util.CheckArguments(cmd, args, 1, 1)
			jsonRpcCallAndPrint("state::force_stop_sub_routine", []interface{}{args[0]})
		} else {
			jsonRpcCallAndPrint("state::kill_sub_routines", args)
		}
	},
}

var autoCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean a stoped auto routine",
	Long: `
clean a stoped auto routine
`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("state::clean_sub_routines", args)
	},
}

func init() {
	autoCmd.Flags().Bool("full-error-checking", false, "Check for errors other than just connectivity errors (default false)")
	autoCmd.Flags().IntP("interval", "i", 50000, "Send interval in microseconds")
	autoCmd.Flags().IntP("send-per-interval", "b", 1, "Send of requests to send per interval tick (default 1)")
	autoKillCmd.Flags().BoolP("force", "f", false, "force kill/stop the routine (this may cause a crash)")
	autoCmd.AddCommand(autoKillCmd, autoCleanCmd)
	RootCmd.AddCommand(autoCmd)
}
