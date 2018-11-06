package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// server int
	testNetID string
)

var testnetCmd = &cobra.Command{
	Use:   "testnet",
	Short: "Get testnet information",
	Long: `Testnet will allow the user to get infromation regarding the test network.

	testnetInfo				Get all testnets which are currently running
	testnetInfo --id [Testnet ID]		Get data on a single testnet
	addTestnet				Add and deploy a new testnet
	`,

	Run: func(cmd *cobra.Command, args []string) {
		curlGET("http://localhost:8000/testnets/" + fmt.Sprint(testNetID))
	},
}

func init() {
	testnetCmd.Flags().StringVarP(&testNetID, "ID", "i", "", "Testnet ID")

	RootCmd.AddCommand(testnetCmd)
}
