package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	testNetID string
)

var testnetCmd = &cobra.Command{
	Use:   "testnet",
	Short: "Get testnet information",
	Long: `Testnet will allow the user to get infromation regarding the test network.

	info				Get all testnets which are currently running
	info id [Testnet ID]		Get data on a single testnet
	`,

	Run: func(cmd *cobra.Command, args []string) {
		curlGET(fmt.Sprint(serverAddr) + "/testnets/" + fmt.Sprint(testNetID))
		// curlGET("http://localhost:8000/testnets/" + fmt.Sprint(testNetID))
	},
}

func init() {
	testnetCmd.Flags().StringVarP(&testNetID, "ID", "i", "", "Testnet ID")
	testnetCmd.Flags().StringVarP(&serverAddr, "serverAddr", "a", "localhost:5000", "server address with port 8000")

	RootCmd.AddCommand(testnetCmd)
}
