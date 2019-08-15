package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"os"
	"runtime"
)

func handleUpdateLinux(branch string) {
	endpoint := fmt.Sprintf("https://storage.cloud.google.com/genesis-public/cli/%s/bin/linux/%s/whiteblock", branch, runtime.GOARCH)
	binary, err := util.HttpRequest("GET", endpoint, "")
	if err != nil {
		util.PrintErrorFatal(err)
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CLI",
	Long:  `Updates the cli binary`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		var branch string
		switch args[0] {
		case "beta":
			fallthrough
		case "dev":
			branch = "dev"
		case "master":
			fallthrough
		case "stable":
			branch = "master"
		default:
			util.PrintErrorFatal("Invalid argument, specify either beta or stable")
		}
		switch runtime.GOOS {
		case "linux":
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
