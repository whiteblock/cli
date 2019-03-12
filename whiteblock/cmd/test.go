package cmd

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/spf13/cobra"
	util "../util"
)

var testCmd = &cobra.Command{
	// Hidden: true,
	Use:    "test <file>",
	Short:  "Run test cases.",
	Long: `

This command will read from a file to run a test.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)

		cwd := os.Getenv("PWD")
		b, err := ioutil.ReadFile(cwd + "/" + args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}

		fmt.Println(prettyp(string(b)))

		var cont map[string]interface{}
		err = json.Unmarshal(b, &cont)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		fmt.Println(cont["build"])
		fmt.Println(cont["netconfig"])
		fmt.Println(cont["rpc"])
		fmt.Println(cont["tests"])

		jsonRpcCallAndPrint("add_commands", cont["rpc"])
		res,err := jsonRpcCall("build", cont["build"])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		buildListener(res.(string))
		jsonRpcCallAndPrint("netem", cont["netconfig"])
		jsonRpcCallAndPrint("run_tests", cont["tests"])
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}

