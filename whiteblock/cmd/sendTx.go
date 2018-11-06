package cmd

import (
	"github.com/spf13/cobra"
)

var (
	txCount     int
	senderCount int
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send transactions from all nodes",
	Long: `Send will have nodes send a specified number of transactions from every node that had been deployed.
	
	transactions			Sends specified number of transactions;
	senders				Number of nodes sending transactions;

	`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	sendCmd.Flags().IntVarP(&txCount, "transactions", "t", 100, "Number of Transactions to Send")
	sendCmd.Flags().IntVarP(&senderCount, "senders", "s", nodes, "Number of Nodes Sending Tx")

	RootCmd.AddCommand(sendCmd)
}
