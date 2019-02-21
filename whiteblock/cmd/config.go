package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Hidden: true,
	Use:    "config",
	Short:  "Configuration file for default parameters for future builds.",
	Run:    PartialCommand,
}

var resetConfCmd = &cobra.Command{
	Hidden: true,
	Use:    "reset",
	Short:  "Reset the config file.",
	Long: `
This command will rest the configuration file when called.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		cwd := os.Getenv("HOME")
		err := os.RemoveAll(cwd + "/.config/whiteblock/config.json")
		if err != nil {
			panic(err)
		}
		fmt.Println("Configuration file has been reset. Run a command to be prompted to create a new one.")
	},
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
			fmt.Println("No configuration file could be found. One will be automatically generated once a successful build has been built. Please refer to the command: 'whiteblock build -h' for help")
			return
		}
		fmt.Println(prettyp(string(b)))
	},
}

func init() {
	confCmd.AddCommand(resetConfCmd, showConfCmd)
	RootCmd.AddCommand(confCmd)
}
