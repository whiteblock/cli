package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
	"strconv"
)

var (
	txsFlag      int
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
	Use:   "tx <command>",
	Short: "Run transaction commands.",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Tx will run commands relevant to sending transactions.

Please use the help commands to make sure you provide the correct flags. If the blockchain is not listed in the help command, the transaction command is not supported for that blockchain. 
	`,
	Run: util.PartialCommand,
}

/*
The primary use of these methods is to be able to send one line commands through a
testing script that will be able to automate transaction tests.
*/

var sendSingleTxCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a transaction between two accounts",
	Long: `
The user must specify the flags that will be used for sending transactions.

Send a transaction between two accounts.

Required Parameters: 
	ethereum:  --from <address>  --destination <address> --gas <gas> --gasprice <gas price> --value <amount>
	eos:  --node <node> --from <address> --destination <address> --value <amount> 
	
Optional Parameters:
	eos:  --symbol [symbol=SYS] --code [code=eosio.token] --memo [memo=]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.RequireFlags(cmd, "from", "destination", "gas", "gasprice", "value")

		util.JsonRpcCallAndPrint("send_transaction", []interface{}{
			util.GetStringFlagValue(cmd, "from"),
			util.GetStringFlagValue(cmd, "destination"),
			util.GetStringFlagValue(cmd, "gas"),
			util.GetStringFlagValue(cmd, "gasprice"),
			strconv.Itoa(util.GetIntFlagValue(cmd, "value")),
		})
	},
}

var sendToTxCmd = &cobra.Command{
	Use:   "to",
	Short: "Send transaction data to an account",
	Long: `
The user must specify the flags that will be used for sending transaction data.
Send a transaction data to an account.

Required Parameters: 
	--destination <address> --value <amount> --data <transaction data>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.RequireFlags(cmd, "destination", "value")

		util.JsonRpcCallAndPrint("send_to", []interface{}{
			util.GetStringFlagValue(cmd, "destination"),
			util.GetStringFlagValue(cmd, "value"),
			util.GetStringFlagValue(cmd, "data"),
		})
	},
}

var startTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "start",
	Short: "Send transactions",
	Long: `
This command will be used to automate transactions and will require require flags to execute. There are two modes of sending transactions: stream and burst. 
The user must specify the flags that will be used for sending transactions.
	`,
	Run: util.PartialCommand,
}

var startStreamTxCmd = &cobra.Command{
	Use:     "stream",
	Short:   "Send continuous transactions",
	Aliases: []string{"cont", "continuous"},
	Long: `
This command will start sending a continual stream of transactions according to the given flags. 
Stream will send transactions as a continuous flow of tps. 
The user will need to run the command tx stop to stop running transactions.
`,

	Run: func(cmd *cobra.Command, args []string) {
		util.RequireFlags(cmd, "tps")

		tps := util.GetIntFlagValue(cmd, "tps")
		newForm, err := cmd.Flags().GetBool("new")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		valueInEth := strconv.Itoa(util.GetIntFlagValue(cmd, "value")) + "000000000000000000"
		size := util.GetIntFlagValue(cmd, "size")

		if !newForm {
			util.JsonRpcCallAndPrint("run_constant_tps", []interface{}{tps, valueInEth, size})
			return
		}

		params := map[string]interface{}{
			"tps":    tps,
			"value":  valueInEth,
			"txSize": size,
			"mode":   util.GetStringFlagValue(cmd, "mode"),
		}

		dest := util.GetStringFlagValue(cmd, "destination")
		if len(dest) > 0 {
			params["destination"] = dest
		}
		log.WithFields(log.Fields{"params": params}).Debug("Sending the request to start sending tx")
		util.JsonRpcCallAndPrint("run_constant_tps", []interface{}{params})
	},
}

var startBurstTxCmd = &cobra.Command{
	Hidden: true,
	Use:    "burst",
	Short:  "Send burst transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
This command will send a burst of transactions. Additional flags are optional. Burst will send one burst of transactions to the blockchain to fill the transaction pool.

Required Parameters: 
	--txs <number of tx>
	--value <value>
Optional Parameters:
	--size [tx size]
`,
	Run: func(cmd *cobra.Command, args []string) {
		params := []string{strconv.Itoa(txsFlag)}
		previousBuild, err := build.GetPreviousBuild()
		if err != nil {
			util.PrintErrorFatal(err)
		}

		switch previousBuild.Blockchain {
		case "eos":
			//error handling for invalid flags
			if valueFlag != 0 {
				util.Print("Invalid \"valueFlag\" flag has been provided.")
				cmd.Help()
				return
			}
			if tpsFlag == 0 {
				util.Print("No \"txsFlag\" flag has been provided. Please input the tps flag with a value.")
				cmd.Help()
				return
			}
			if txSizeFlag >= 174 {
				params = append(params, strconv.Itoa(txSizeFlag))
			} else if txSizeFlag > 0 && txSizeFlag < 174 {
				util.Print("Transaction size value is too small. The minimum size of a transaction is 174 bytes.")
				return
			}
		default:
			util.ClientNotSupported(previousBuild.Blockchain)
		}
		util.JsonRpcCallAndPrint("run_burst_tx", params)
	},
}

var stopTxCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop transactions",
	Long: `
The user must specify the blockchain flag as well as any other flags that will be used for sending transactions.
Stops the sending of transactions if transactions are currently being sent
`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := util.JsonRpcCall("state::kill", []string{})
		if res != nil && res.(float64) == 0 && err == nil {
			util.Print("Transactions stopped successfully")
		} else {
			util.PrintErrorFatal("There was an error stopping transactions")
		}
	},
}

func init() {

	sendSingleTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	sendSingleTxCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "where the transaction will be sent from")
	sendSingleTxCmd.Flags().StringVarP(&gasFlag, "gas", "g", "", "specify gas for tx")
	sendSingleTxCmd.Flags().StringVarP(&nodeFlag, "node", "n", "", "specify node to send tx")
	sendSingleTxCmd.Flags().StringVarP(&gasPriceFlag, "gasprice", "p", "", "specify gas price for tx")
	sendSingleTxCmd.Flags().IntVarP(&valueFlag, "value", "v", 0, "amount to send in transaction")

	sendToTxCmd.Flags().StringP("destination", "d", "", "where the transaction will be sent to")
	sendToTxCmd.Flags().StringP("value", "v", "", "amount to send in transaction")
	sendToTxCmd.Flags().String("data", "0x0", "transaction data")

	startStreamTxCmd.Flags().StringP("destination", "d", "", "where the transaction will be sent to")
	startStreamTxCmd.Flags().IntP("size", "s", -1, "size of the transaction in bytes")
	startStreamTxCmd.Flags().IntP("tps", "t", 0, "transactions per second")
	startStreamTxCmd.Flags().IntP("value", "v", -1, "amount to send in transaction")
	startStreamTxCmd.Flags().String("mode", "", "the tx send mode: nr,default")
	startStreamTxCmd.Flags().Bool("new", false, "use the new tx format")
	startStreamTxCmd.MarkFlagRequired("tps")

	startBurstTxCmd.Flags().StringVarP(&toFlag, "destination", "d", "", "where the transaction will be sent to")
	startBurstTxCmd.Flags().IntVarP(&txSizeFlag, "size", "s", 0, "size of the transaction in bytes")
	startBurstTxCmd.Flags().IntVarP(&txsFlag, "txs", "t", 0, "transactions per second")
	startBurstTxCmd.Flags().IntVarP(&valueFlag, "value", "v", -1, "amount to send in transaction")

	startTxCmd.AddCommand(startStreamTxCmd, startBurstTxCmd)
	txCmd.AddCommand(sendSingleTxCmd, startTxCmd, stopTxCmd)
	sendSingleTxCmd.AddCommand(sendToTxCmd)
	RootCmd.AddCommand(txCmd)
}
