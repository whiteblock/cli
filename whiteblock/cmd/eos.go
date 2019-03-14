package cmd

import (
	util "../util"
	"github.com/spf13/cobra"
)

var eosCmd = &cobra.Command{
	Use:   "eos <command>",
	Short: "Run eos commands",
	Long:  "\nEos will allow the user to get information and run EOS commands.\n",
	Run:   util.PartialCommand,
}

var eosGetInfoCmd = &cobra.Command{
	Use:   "get_info [node]",
	Short: "Get EOS info",
	Long: `
Roughly equivalent to calling cleos get info

Params: The node to get info from
Format: [node]

Response: eos blockchain state info`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 0, 1)
		jsonRpcCallAndPrint("eos::get_info", args)
	},
}

/*

var eosSendTxCmd = &cobra.Command{
	Use:   "send_transaction <node> <from> <to> <amount> [symbol=SYS] [code=eosio.token] [memo=]",
	Short: "Send single transaction to another account",
	Long: `
This command will send a single transaction from one account to another. Additional arguments are required.

Params: node number, from account, to account, amount, symbol, code, memo
Format: <node> <from> <to> <amount> [symbol=SYS] [code=eosio.token] [memo=]

Response: The txid`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args,4,7)
		jsonRpcCallAndPrint("eos::send_transaction",args)
	},
}

var eosSendBurstTxCmd = &cobra.Command{
	Use:   "send_burst_transaction <tps> [tx size]",
	Short: "Send burst transactions",
	Long: `
This command will send a burst of transactions. Additional arguments are optional.

Params: number of transactions to send per second, transaction size
Format: <txs>, [tx size]

Response: success or ERROR`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args,1,2)
		jsonRpcCallAndPrint("eos::run_burst_tx",args)
	},
}

var eosConstTpsCmd = &cobra.Command{
	Use:   "run_constant_tps <tps> [tx size]",
	Short: "Send continuous transactions",
	Long: `
This command will have all nodes send continous transactions at a constant rate.

Params: number of transactions to send per second, transaction size
Format: <tps>, [tx size]

Response: success or ERROR`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args,1,2)
		tps, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidInteger("tps",args[0],true)
		}
		if tps > 5000 {
			fmt.Println("The limit for tps is set to 5000. Please input a lower value.")
			return
		}
		jsonRpcCallAndPrint("eos::run_constant_tps",args)
	},
}

*/

func init() {
	//eos subcommands
	eosCmd.AddCommand(eosGetInfoCmd)

	RootCmd.AddCommand(eosCmd)
}
