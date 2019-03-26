package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

type Server struct {
	Addr     string   `json:"addr"`  //IP to access the server
	Nodes    int      `json:"nodes"`
	Max      int      `json:"max"`
	Id       int      `json:"id"`
	SubnetID int      `json:"subnetID"`
	Ips      []string `json:"ips"`
}

var RootCmd = &cobra.Command{
	Use:   "whiteblock",
	Short: "Create and test blockchains",
	Long: `This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation, usages, and exmaples can be found at www.whiteblock.io/docs/cli.
	`,
}
func init(){
	RootCmd.PersistentFlags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
}
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}