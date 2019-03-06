package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	util "../util"
)

type List struct {
	Results []struct {
		Series []struct {
			Columns []string        `json:"columns"`
			Values  [][]interface{} `json:"values"`
		}
	} `json:"results"`
}

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

var testStartCMD = &cobra.Command{
	Use:   "start <minimum latency> <minimum completion percentage> <number of assets to send> <asset sends per block>",
	Short: "Starts propagation test.",
	Long: `
Sys test start will start the propagation test. It will wait for the signal start time, have nodes send messages at the same time, and require to wait a minimum amount of time then check receivers with a completion rate of minimum completion percentage. 

Format: <minimum latency> <minimum completion percentage> <number of assets to send> <asset sends per block>
Params: Time in seconds, percentage, number of assets to send, asset sends per block

`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 4, 4)
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
		util.CheckArguments(args, 1, 1)
		results, err := jsonRpcCall("sys::get_recent_test_results", args)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		result, ok := results.(map[string]interface{})
		if !ok {
			panic(1)
		}

		rc := result["results"].([]interface{})[0]
		s := rc.(map[string]interface{})["series"]
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
	testStartCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	// testResultsCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	sysTestCMD.AddCommand(testStartCMD, testResultsCMD)
	sysCMD.AddCommand(sysTestCMD)
	RootCmd.AddCommand(sysCMD)
}
