package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	tpsFlag      int
	valueFlag    int
	txSizeFlag   int
	gasFlag      string
	gasPriceFlag string
	toFlag       string
	fromFlag     string
	nodeFlag     string
)

var txCmd = &cobra.Command{
	Hidden: true,
	Use:    "tx <command>",
	Short:  "Run transaction commands.",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Tx will run commands relavent to sending transactions.

Please use the help commands to make sure you provide the correct flags. If the blockchain is not listed in the help command, the transaction command is not supported for that blockchain. 
	`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list above.\n")
		cmd.Help()
		return
	},
}

/*
We have to figure out a generalized format of how the user will enter the arguments.
Right now, geth, eos, and sys all have different parameters they need to provide to
send tx. I am just using the geth description of the send_tx as a place holder for
now. This is for all the commands in the file.

The primary use of these methods is to be able to send one line commands through a
testing script that will be able to automate transaction tests.
*/

var startGeneralizedTxCmd = &cobra.Command{
	Hidden:  true,
	Use:     "send_const_tx",
	Aliases: []string{"servers"},
	Short:   "Send continuous transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
This command will start sending a continual stream of transactions according to the given flags, --value = -1 means randomize value.

Required Parameters: 
	ethereum:  --tps <tps> --value <value>
	eos:  --tps <tps> 

Optional Parameters:
	ethereum:  --to <address>
	eos:  --size <tx size>

	`,

	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(command)
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		if blockchain != blockchainFlag {
			fmt.Println("Wrong blockchain. You are connected to the " + blockchain + " blockchain")
			return
		}
		switch blockchainFlag {
		case "":
			fmt.Println("No blockchain has been provided. Please indicate the type of blockchain using the flag.")
		case "ethereum":
			command = "eth::start_transactions"
			toEth := strconv.Itoa(valueFlag) + "000000000000000000"
			param = strconv.Itoa(tpsFlag) + " " + toEth
			if len(toFlag) > 0 {
				param = param + " " + toFlag
			}
		case "eos":
			command = "eos::run_constant_tps"
			param = strconv.Itoa(tpsFlag) + " " + strconv.Itoa(txSizeFlag)
		case "syscoin":
			command = "sys::start_test"
			// I think we need to change how the test will be sent in the backend if we want to generalize the transactions for syscoin
			// param = "{\"waitTime\":" + args[0] + ",\"minCompletePercent\":" + args[1] + ",\"numberOfTransactions\":" + args[2] + "}"
		}
		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var startGeneralizedBurstTxCmd = &cobra.Command{
	Hidden: true,
	Use:    "send_burst_tx",
	Short:  "Send burst transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
This command will send a burst of transactions. Additional flags are optional.

Required Parameters: 
	eos:  --tps <tps>  --size <tx size>


`,
	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		switch blockchainFlag {
		case "":
			fmt.Println("No blockchain has been provided. Please indicate the type of blockchain using the flag.")
		case "ethereum":
			fmt.Println("This function is not supported for the ethereum client.")
		case "eos":
			command = "eos::run_burst_tx"
			param = strconv.Itoa(tpsFlag) + " " + strconv.Itoa(txSizeFlag)
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
		}

		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var sendGeneralizedTxCmd = &cobra.Command{
	Hidden: true,
	Use:    "send_tx",
	Short:  "Sends a single transaction",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Send a transaction between two accounts. 

Required Parameters: 
	ethereum:  --from <address>  --to <address> --gas <gas> --gasprice <gas price> --value <value to send>
	eos:  --node <node number> --from <address> --to <address> --value <value to send> 
	
Optional Parameters:
	eos:  --symbol [symbol=SYS] --code [code=eosio.token] --memo [memo=]

`,
	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		switch blockchainFlag {
		case "":
			fmt.Println("No blockchain has been provided. Please indicate the type of blockchain using the flag.")
		case "ethereum":
			command = "eth::send_transaction"
			param = fromFlag + " " + toFlag + " " + gasFlag + " " + gasPriceFlag + " " + strconv.Itoa(valueFlag)
		case "eos":
			command = "eos::send_transaction"
			param = nodeFlag + " " + fromFlag + " " + toFlag + " " + strconv.Itoa(valueFlag)
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var stopGeneralizedTxCmd = &cobra.Command{
	Hidden: true,
	Use:    "stop_tx",
	Short:  "Stop transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Stops the sending of transactions if transactions are currently being sent
`,
	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		switch blockchainFlag {
		case "":
			fmt.Println("No blockchain has been provided. Please indicate the type of blockchain using the flag.")
		case "ethereum":
			command = "eth::stop_transactions"
		case "eos":
			command = "eth::stop_transactions"
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
		}

		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

func init() {
	txCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	startGeneralizedTxCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	startGeneralizedBurstTxCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	startGeneralizedTxCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	startGeneralizedTxCmd.Flags().StringVarP(&toFlag, "to", "d", "", "where the transaction will be sent to")
	startGeneralizedTxCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "where the transaction will be sent from")
	startGeneralizedTxCmd.Flags().IntVarP(&txSizeFlag, "size", "s", 0, "size of the transaction in bytes")
	startGeneralizedTxCmd.Flags().IntVarP(&tpsFlag, "tps", "t", 0, "transactions per second")
	startGeneralizedTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	startGeneralizedBurstTxCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	startGeneralizedBurstTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	startGeneralizedBurstTxCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "where the transaction will be sent from")
	startGeneralizedBurstTxCmd.Flags().IntVarP(&txSizeFlag, "size", "s", 0, "size of the transaction in bytes")
	startGeneralizedBurstTxCmd.Flags().IntVarP(&tpsFlag, "tps", "t", 0, "transactions per second")
	startGeneralizedBurstTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	sendGeneralizedTxCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	sendGeneralizedTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	sendGeneralizedTxCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "where the transaction will be sent from")
	sendGeneralizedTxCmd.Flags().StringVarP(&gasFlag, "gas", "g", "", "specify gas for tx")
	sendGeneralizedTxCmd.Flags().StringVarP(&nodeFlag, "node", "n", "", "specify node to send tx")
	sendGeneralizedTxCmd.Flags().StringVarP(&gasPriceFlag, "gasprice", "p", "", "specify gas price for tx")
	sendGeneralizedTxCmd.Flags().IntVarP(&txSizeFlag, "size", "s", 0, "size of the transaction in bytes")
	sendGeneralizedTxCmd.Flags().IntVarP(&tpsFlag, "tps", "t", 0, "transactions per second")
	sendGeneralizedTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	stopGeneralizedTxCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")

	txCmd.AddCommand(startGeneralizedTxCmd, startGeneralizedBurstTxCmd, sendGeneralizedTxCmd, stopGeneralizedTxCmd)
	RootCmd.AddCommand(txCmd)
}
