package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	util "../util"
)


var sysCMD = &cobra.Command{
	Use:   "sys <command>",
	Short: "Run SYS commands.",
	Long:  "\nSys will allow the user to get information and run SYS commands.\n",
	Run:   util.PartialCommand,
}

var sysTestCMD = &cobra.Command{
	Use:   "test <command>",
	Short: "SYS test commands.",
	Long:  "\nSys test will allow the user to get infromation and run SYS tests.\n",
	Run:   util.PartialCommand,
}

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
Format: [node]

Response: eos blockchain state info`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd,args, 0, 1)
		jsonRpcCallAndPrint("eos::get_info", args)
	},
}



var testStartCMD = &cobra.Command{
	Use:   "start <minimum latency> <minimum completion percentage> <number of assets to send> <asset sends per block>",
	Short: "Starts propagation test.",
	Long: `
Sys test start will start the propagation test. It will wait for the signal start time, have nodes send messages at the same time, and require to wait a minimum amount of time then check receivers with a completion rate of minimum completion percentage. 

Format: <minimum latency> <minimum completion percentage> <number of assets to send> <asset sends per block>
Params: Time in seconds, percentage, number of assets to send, asset sends per block

`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd,args, 4, 4)
		jsonRpcCallAndPrint("sys::start_test", args)
		return
	},
}

var testResultsCMD = &cobra.Command{
	Use:   "results <test number>",
	Short: "Get results from a previous test.",
	Long: `
Sys test results pulls data from a previous test or tests and outputs as csv.

Format: <test number>
Params: Test number

	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd,args, 1, 1)
		results, err := jsonRpcCall("sys::get_recent_test_results", args)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		result, ok := results.(map[string]interface{})
		if !ok {
			util.PrintStringError("Got back the results in an invalid format")
		}

		rc := result["results"].([]interface{})[0]
		s := rc.(map[string]interface{})["series"]
		if s == nil{
			fmt.Println("No results availible")
			os.Exit(1)
		}
		sc := s.([]interface{})[0].(map[string]interface{})
		c := sc["columns"].([]interface{})
		v := sc["values"].([]interface{})

		for i := 0; i < len(v); i++ {
			fmt.Printf("[%d]\n", i)
			vv := v[i].([]interface{})
			for j := 0; j < len(vv); j++ {
				fmt.Printf("\t%v: %v\n", c[j], vv[j])
			}
		}
	},
}

func init() {
	eosCmd.AddCommand(eosGetInfoCmd)

	sysTestCMD.AddCommand(testStartCMD, testResultsCMD)
	sysCMD.AddCommand(sysTestCMD)

	RootCmd.AddCommand(sysCMD,eosCmd)
}
