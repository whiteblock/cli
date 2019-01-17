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
	// Hidden: true,
	Use:   "tx <command>",
	Short: "Run transaction commands.",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Tx will run commands relavent to sending transactions.

Please use the help commands to make sure you provide the correct flags. If the blockchain is not listed in the help command, the transaction command is not supported for that blockchain. 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.")
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

var sendSingleTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "send",
	Short: "Sends a single transaction",
	Long: `
The user must specify the flags that will be used for sending transactions.
Send a transaction between two accounts.

Required Parameters: 
	ethereum:  --from <address>  --destination <address> --gas <gas> --gasprice <gas price> --value <amount>
	eos:  --node <node number> --from <address> --destination <address> --value <amount> 
	
Optional Parameters:
	eos:  --symbol [symbol=SYS] --code [code=eosio.token] --memo [memo=]

`,
	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		switch blockchain {
		case "ethereum":
			if !(len(toFlag) > 0) || !(len(fromFlag) > 0) || !(len(gasFlag) > 0) || !(len(gasPriceFlag) > 0) || valueFlag == 0 {
				fmt.Println("Required flags were not provided. Please input the required flags.")
				return
			}
			command = "eth::send_transaction"
			param = fromFlag + " " + toFlag + " " + gasFlag + " " + gasPriceFlag + " " + strconv.Itoa(valueFlag)
		case "eos":
			if !(len(nodeFlag) > 0) || !(len(toFlag) > 0) || !(len(fromFlag) > 0) || valueFlag == 0 {
				fmt.Println("Required flags were not provided. Please input the required flags.")
				return
			}
			command = "eos::send_transaction"
			param = nodeFlag + " " + fromFlag + " " + toFlag + " " + strconv.Itoa(valueFlag)
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			println(blockchain)
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var startTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "start",
	Short: "Send transactions",
	Long: `
The user must specify the flags that will be used for sending transactions.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.")
		cmd.Help()
		return
	},
}

var startStreamTxCmd = &cobra.Command{
	// Hidden: true,
	Use:     "stream",
	Short:   "Send continuous transactions",
	Aliases: []string{"cont", "continuous"},
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
This command will start sending a continual stream of transactions according to the given flags, --value = -1 means randomize value.

Required Parameters: 
	ethereum:  --tps <tps> --value <amount>
	eos:  --tps <tps> 

Optional Parameters:
	ethereum:  --destination [address]
	eos:  --size [tx size]
	`,

	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		switch blockchain {
		case "ethereum":
			//error handling for invalid flags
			if !(txSizeFlag == 0) {
				fmt.Println("Invalid use of flag \"txSizeFlag\". This is not supported with Ethereum")
				return
			}
			if valueFlag == 0 {
				fmt.Println("No \"valueFlag\" has been provided. Please input the value flag with a value.")
				return
			}
			if tpsFlag == 0 {
				fmt.Println("No \"tpsFlag\" flag has been provided. Please input the tps flag with a value.")
				return
			}

			command = "eth::start_transactions"
			toEth := strconv.Itoa(valueFlag) + "000000000000000000"
			param = strconv.Itoa(tpsFlag) + " " + toEth
			if len(toFlag) > 0 {
				param = param + " " + toFlag
			}
		case "eos":
			//error handling for invalid flags
			if valueFlag != 0 {
				fmt.Println("Invalid \"valueFlag\" flag has been provided.")
				return
			}
			if tpsFlag == 0 {
				fmt.Println("No \"tpsFlag\" flag has been provided. Please input the tps flag with a value.")
				return
			}

			command = "eos::run_constant_tps"
			param = strconv.Itoa(tpsFlag)
			if txSizeFlag >= 174 {
				param = param + " " + strconv.Itoa(txSizeFlag)
			} else if txSizeFlag > 0 && txSizeFlag < 174 {
				fmt.Println("Transaction size value is too small. The minimum size of a transaction is 174 bytes.")
				return
			}
		case "syscoin":
			command = "sys::start_test"
			return
			// I think we need to change how the test will be sent in the backend if we want to generalize the transactions for syscoin
			// param = "{\"waitTime\":" + args[0] + ",\"minCompletePercent\":" + args[1] + ",\"numberOfTransactions\":" + args[2] + "}"
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var startBurstTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "burst",
	Short: "Send burst transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
This command will send a burst of transactions. Additional flags are optional.

Required Parameters: 
	eos:  --tps <tps>  
Optional Parameters:
	--size [tx size]
`,
	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		switch blockchain {
		case "ethereum":
			fmt.Println("This function is not supported for the ethereum client.")
			return
		case "eos":
			//error handling for invalid flags
			if valueFlag != 0 {
				fmt.Println("Invalid \"valueFlag\" flag has been provided.")
				return
			}
			if tpsFlag == 0 {
				fmt.Println("No \"tpsFlag\" flag has been provided. Please input the tps flag with a value.")
				return
			}

			command = "eos::run_burst_tx"
			param = strconv.Itoa(tpsFlag)
			if txSizeFlag >= 174 {
				param = param + " " + strconv.Itoa(txSizeFlag)
			} else if txSizeFlag > 0 && txSizeFlag < 174 {
				fmt.Println("Transaction size value is too small. The minimum size of a transaction is 174 bytes.")
				return
			}
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var stopTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "stop",
	Short: "Stop transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Stops the sending of transactions if transactions are currently being sent
`,
	Run: func(cmd *cobra.Command, args []string) {
		command := ""
		param := ""
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		switch blockchain {
		case "ethereum":
			command = "eth::stop_transactions"
		case "eos":
			command = "eth::stop_transactions"
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		fmt.Println("Stopped transactions.")
		wsEmitListen(serverAddr, command, param)
	},
}

func init() {
	txCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	sendSingleTxCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	sendSingleTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	sendSingleTxCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "where the transaction will be sent from")
	sendSingleTxCmd.Flags().StringVarP(&gasFlag, "gas", "g", "", "specify gas for tx")
	sendSingleTxCmd.Flags().StringVarP(&nodeFlag, "node", "n", "", "specify node to send tx")
	sendSingleTxCmd.Flags().StringVarP(&gasPriceFlag, "gasprice", "p", "", "specify gas price for tx")
	sendSingleTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	startStreamTxCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	startStreamTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	startStreamTxCmd.Flags().IntVarP(&txSizeFlag, "size", "s", 0, "size of the transaction in bytes")
	startStreamTxCmd.Flags().IntVarP(&tpsFlag, "tps", "t", 0, "transactions per second")
	startStreamTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	startBurstTxCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	startBurstTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	startBurstTxCmd.Flags().IntVarP(&txSizeFlag, "size", "s", 0, "size of the transaction in bytes")
	startBurstTxCmd.Flags().IntVarP(&tpsFlag, "tps", "t", 0, "transactions per second")
	startBurstTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	stopTxCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	startTxCmd.AddCommand(startStreamTxCmd, startBurstTxCmd)
	txCmd.AddCommand(sendSingleTxCmd, startTxCmd, stopTxCmd)
	RootCmd.AddCommand(txCmd)
}
