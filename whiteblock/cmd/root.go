package cmd

import (
	"os"
	"fmt"
	"encoding/json"
	"github.com/spf13/cobra"
	util "../util"
)

/*
	Globals
 */

var (
	serverAddr	string
)

var RootCmd = &cobra.Command{
	Use:   "whiteblock",
	Short: "Create and test blockchains",
	Long: `This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation, usages, and exmaples can be found at www.whiteblock.io/docs/cli.
	`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


var completionCmd = &cobra.Command{
	Hidden:true,
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(whiteblock completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(whiteblock completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout);
	},
}

func init(){
	RootCmd.PersistentFlags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	RootCmd.AddCommand(completionCmd)
	//Possibly update this on load.
	if util.StoreExists("profile") {
		rawProfile,err := util.ReadStore("profile")
		if err != nil{
			panic(err)
		}
		var profile map[string]interface{}
		err = json.Unmarshal(rawProfile,&profile)
		if err != nil {
			panic(err)
		}
		biomes,ok := profile["biomes"].([]interface{})
		if !ok {
			//If there aren't any biomes for the jwt, don't continue or try fetching?
			return
		}
		if len(biomes) == 0 {
			return
		}

		biome,ok := biomes[0].(map[string]interface{})
		if !ok {
			return
		}
		host,ok := biome["host"].(string)
		if !ok {
			return
		}
		serverAddr = host + ":5001"
		
	}
}