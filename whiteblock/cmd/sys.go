package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var sysCMD = &cobra.Command{
	Use:   "sys <command>",
	Short: "Run SYS commands.",
	Long: `
Sys will allow the user to get infromation and run SYS commands.
`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock sys -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var sysTestCMD = &cobra.Command{
	Use:   "test <command>",
	Short: "SYS test commands.",
	Long: `
Sys test will allow the user to get infromation and run SYS tests.
`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock sys test -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var testStartCMD = &cobra.Command{
	Use:   "start <wait time> <min complete percent> <number of tx>",
	Short: "Starts propagation test.",
	Long: `
Sys test start will start the propagation test. It will wait for the signal start time, have nodes send messages at the same time, and require to wait a minimum amount of time then check receivers with a completion rate of minimum completion percentage. 

Format: <wait time> <min complete percent> <number of tx>
Params: Time in seconds, percentage, number of transactions

`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "sys::start_test"
		if len(args) != 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock sys test start -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}
		param := "{\"waitTime\":" + args[0] + ",\"minCompletePercent\":" + args[1] + ",\"numberOfTransactions\":" + args[2] + "}"
		// println(command)
		// println(param)
		wsEmitListen(serverAddr, command, param)
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "sys::get_recent_test_results"
		if len(args) != 1 {
			out, err := exec.Command("bash", "-c", "./whiteblock sys test results -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}
		wsEmitListen(serverAddr, command, args[0])
	},
}

func init() {
	testStartCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	testResultsCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	sysTestCMD.AddCommand(testStartCMD, testResultsCMD)
	sysCMD.AddCommand(sysTestCMD)
	RootCmd.AddCommand(sysCMD)
}
