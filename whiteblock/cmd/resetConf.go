package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var resetConfCmd = &cobra.Command{
	Use:   "reset-conf",
	Short: "Reset the configuration file.",
	Long: `

This command will rest the configuration file when called.
	`,

	Run: func(cmd *cobra.Command, args []string) {

		cwd := os.Getenv("HOME")
		err := os.RemoveAll(cwd + "/cli/whiteblock/config")
		if err != nil {
			panic(err)
		}
		println("Configuration file has been reset. Run a command to be prompted to create a new one.")
	},
}

func init() {
	RootCmd.AddCommand(resetConfCmd)
}
