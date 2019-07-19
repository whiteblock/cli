package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
)

var confCmd = &cobra.Command{
	Hidden: true,
	Use:    "config",
	Short:  "Configuration file for default parameters for future builds.",
	Run:    util.PartialCommand,
}

var showConfCmd = &cobra.Command{
	Hidden: true,
	Use:    "show",
	Short:  "Show config file",
	Long: `
	Show the default values set by the configuration file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.Print(*conf)
	},
}

func init() {
	confCmd.AddCommand(showConfCmd)
	RootCmd.AddCommand(confCmd)
}
