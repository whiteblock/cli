package cmd

import (
	"github.com/spf13/cobra"
)

var (
	gethcommand string
)

var gethCmd = &cobra.Command{
	Use:   "geth [command]",
	Short: "Run geth commands",
	Long: `Geth will allow the user to get infromation and run geth commands.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if gethcommand == "" {

		}
	},
}

var getBlockNumberCmd = &cobra.Command{
	Use:   "get_block_number",
	Short: "Get block number",
	Long: `Get the current highest block number of the chain
	Response: The block number e.g. 10`,
}

var getBlockCmd = &cobra.Command{
	Use:   "get_block",
	Short: "Get block information",
	Long: `Get the data of a block
	Response: JSON Representation of the block.
	
	Params: Block number
	Format: <Block Number>`,
}

var getAccountCmd = &cobra.Command{
	Use:   "get_accounts",
	Short: "Get account information",
	Long: `Get a list of all unlocked accounts
	Response: A JSON array of the accounts`,
}

var getBalanceCmd = &cobra.Command{
	Use:   "get_balance [address]",
	Short: "Get account balance information",
	Long: `Get the current balance of an account
	Response: The integer balance of the account in wei
	
	Params: Account address
	Format: <address>`,
}

var sendTxCmd = &cobra.Command{
	Use:   "send_transaction",
	Short: "Sends a transaction",
	Long: `Send a transaction between two accounts
	Response: The transaction hash

	Params: Sending account, receiving account, gas, gas price, amount to send, transaction data, nonce
	Format: <from> <to> <gas> <gas price> <value> [data] [nonce]
	`,
}

var getTxCountCmd = &cobra.Command{
	Use:   "get_transaction_count",
	Short: "Get transaction count",
	Long: `Get the transaction count sent from an address, optionally by block
	Response: The transaction count
	
	Params: The sender account, a block number
	Format: <address> [block number]`,
}

var getTxCmd = &cobra.Command{
	Use:   "get_transaction",
	Short: "Get transaction information",
	Long: `Get a transaction by its hash
	Response: JSON representation of the transaction.
	
	Params: The transaction hash
	Format: <hash>`,
}

var getTxReceiptCmd = &cobra.Command{
	Use:   "get_transaction_receipt",
	Short: "Get transaction receipt",
	Long: `Get the transaction receipt by the tx hash
	Response: JSON representation of the transaction receipt.
	
	Params: The transaction hash
	Format: <hash>`,
}

var getHashRateCmd = &cobra.Command{
	Use:   "get_hash_rate",
	Short: "Get hasg rate",
	Long: `Get the current hash rate per node
	Response: The hash rate of a single node in the network`,
}

var startTxCmd = &cobra.Command{
	Use:   "start_transactions",
	Short: "Start transactions",
	Long: `Start sending transactions according to the given parameters, value = -1 means randomize value.
	
	Params: The amount of transactions to send in a second, the value of each transaction in wei, the destination for the transaction
	Format: <tx/s> <value> [destination]`,
}

var stopTxCmd = &cobra.Command{
	Use:   "stop_transactions",
	Short: "Start transactions",
	Long:  `Stops the sending of transactions if transactions are currently being sent`,
}

var startMiningCmd = &cobra.Command{
	Use:   "start_mining",
	Short: "Start Mining",
	Long: `Send the start mining signal to nodes, may take a while to take effect due to DAG generation
	Response: The number of nodes which successfully received the signal to start mining
	
	Params: A list of the nodes to start mining or None for all nodes
	Format: [node 1 number] [node 2 number]...`,
}

var stopMiningCmd = &cobra.Command{
	Use:   "stop_mining",
	Short: "Stop mining",
	Long: `Send the stop mining signal to nodes
	Response: The number of nodes which successfully received the signal to stop mining
	
	Params: A list of the nodes to stop mining or None for all nodes
	Format: [node 1 number] [node 2 number]...`,
}

var blockListenerCmd = &cobra.Command{
	Use:   "block_listener",
	Short: "Get block listener",
	Long: `Get all blocks and continue to subscribe to new blocks
	Response: Will emit on eth::block_listener for every block after the given block or 0 that exists/has been created
	
	Params: The block number to start at or None for all blocks
	Format: [block number]`,
}

var getRecentSentTxCmd = &cobra.Command{
	Use:   "get_recent_sent_tx",
	Short: "Get recently sent transaction",
	Long: `Get a number of the most recent transactions sent
	
	Response: JSON object of transaction data
	
	Params: The number of transactions to retrieve
	Format: [number]`,
}

func init() {
	gethCmd.Flags().StringVarP(&gethcommand, "command", "c", "", "Geth command")
	gethCmd.Flags().StringVarP(&serverAddr, "serverAddr", "a", "localhost:5000", "server address with port 8000")

	//geth subcommands

	RootCmd.AddCommand(gethCmd)
}
