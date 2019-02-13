package cmd

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	// Hidden: true,
	Use:    "test <file>",
	Short:  "Run test cases.",
	Long: `

This command will read from a file to run a test.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckArguments(args, 1, 1)

		cwd := os.Getenv("HOME")
		b, err := ioutil.ReadFile(cwd + "/tests/" + args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println(prettyp(string(b)))

		var cont map[string]interface{}
		err = json.Unmarshal(b, &cont)
		if err != nil {
			panic(err)
		}

		fmt.Println(cont["build"])
		fmt.Println(cont["netconfig"])
		fmt.Println(cont["rpc"])
		fmt.Println(cont["test"])

		jsonRpcCallAndPrint("add_commands", cont["rpc"])
		jsonRpcCallAndPrint("build", cont["build"])
		buildListener()
		jsonRpcCallAndPrint("netem", cont["netconfig"])
		jsonRpcCallAndPrint("run_tests", cont["tests"])
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}

