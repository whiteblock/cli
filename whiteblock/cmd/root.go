package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"os"
)

/*
	Globals
*/

var (
	serverAddr string
	conf       = util.GetConfig()
)

var RootCmd = &cobra.Command{
	Use:     "whiteblock",
	Version: VERSION,
	Short:   "Create and test blockchains",
	Long: `This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network.
	Documentation, usages, and exmaples can be found at https://docs.whiteblock.io/.
	`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var completionCmd = &cobra.Command{
	Hidden: true,
	Use:    "completion",
	Short:  "Generates bash completion scripts",
	Long: `To load completion run
. <(whiteblock completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	RootCmd.AddCommand(completionCmd)
	//Possibly update this on load.
	if util.StoreExists("profile") {

		err := LoadProfile() //Load the profile into the profile global
		if err != nil {
			util.DeleteStore("profile")
			util.PrintErrorFatal(err)
		}

		err = LoadBiomeAddress()
		if err != nil {
			util.PrintErrorFatal(err)
		}
	}
}
