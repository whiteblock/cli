package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os"
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

		cwd := os.Getenv("HOME")
		b, err := ioutil.ReadFile(cwd + "/.config/whiteblock/config.json")
		if err != nil {
			util.PrintErrorFatal("No configuration file could be found. One will be automatically generated once a successful build has been built. Please refer to the command: 'whiteblock build -h' for help")
			return
		}
		util.Print(string(b))
	},
}

func init() {
	confCmd.AddCommand(showConfCmd)
	RootCmd.AddCommand(confCmd)
}
