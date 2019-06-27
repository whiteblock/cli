package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os"
	"sort"
)

func GetNodes() ([]Node, error) {
	testnetId, err := getPreviousBuildId()
	if err != nil {
		return nil, err
	}
	res, err := util.JsonRpcCall("nodes", []string{testnetId})
	if err != nil {
		return nil, err
	}
	tmp, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	var out []Node
	log.WithFields(log.Fields{"res": res}).Trace("raw nodes")
	return out, json.Unmarshal(tmp, &out)
}

func readContractsFile() ([]byte, error) {
	cwd := os.Getenv("HOME")
	return ioutil.ReadFile(cwd + "/smart-contracts/whiteblock/contracts.json")
}

var getCmd = &cobra.Command{
	Use:   "get <command>",
	Short: "Get server and network information.",
	Long:  "\nGet will output server and network information and statistics.\n",
	Run:   util.PartialCommand,
}

var getServerCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"servers"},
	Short:   "Get server information.",
	Long:    "\nServer will output server information.\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("get_servers", []string{})
	},
}

var getTestnetIDCmd = &cobra.Command{
	Use:     "testnetid",
	Aliases: []string{"id"},
	Short:   "Get the last stored testnet id",
	Long:    "\nGet the last stored testnet id.\n",
	Run: func(cmd *cobra.Command, args []string) {
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		fmt.Println(testnetID)
	},
}

var getBuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"built"},
	Short:   "Get the last applied build",
	Long:    "\nGet the last applied build.\n",
	Run: func(cmd *cobra.Command, args []string) {
		prevBuild, err := getPreviousBuild()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		fmt.Println(util.Prettypi(prevBuild))
	},
}

var getSupportedCmd = &cobra.Command{
	Use:     "supported",
	Aliases: []string{"blockchains"},
	Short:   "Get the currently supported blockchains",
	Long:    "Fetches the blockchains which whiteblock is currently able build by default",
	Run: func(cmd *cobra.Command, args []string) {

		var blockchains []string
		util.JsonRpcCallP("get_supported_blockchains", []string{}, &blockchains)
		sortedBlockchains := sort.StringSlice(blockchains)
		sortedBlockchains.Sort()
		fmt.Println(util.Prettypi([]string(sortedBlockchains)))
	},
}

var getNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Nodes will show all nodes in the network.",
	Long:    "\nNodes will output all of the nodes in the current network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if all {
			util.JsonRpcCallAndPrint("status_nodes", []string{testnetID})
			return
		}
		res, err := util.JsonRpcCall("status_nodes", []string{testnetID})
		if err != nil {
			util.PrintErrorFatal(err)
		}

		rawNodes := res.([]interface{})
		out := []interface{}{}
		for _, rawNode := range rawNodes {
			if rawNode.(map[string]interface{})["up"].(bool) {
				out = append(out, rawNode)
			}
		}
		fmt.Println(util.Prettypi(out))
	},
}

var getRunningCmd = &cobra.Command{
	Use:   "running",
	Short: "Running will check if a test is running.",
	Long: `
Running will check whether or not there is a test running and get the name of the currently running test.

Response: true or false, on whether or not a test is running; The name of the test or nothing if there is not a test running.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		util.JsonRpcCallAndPrint("state::is_running", []string{})
		util.JsonRpcCallAndPrint("state::what_is_running", []string{})
	},
}

var getDefaultsCmd = &cobra.Command{
	Use:     "default <blockchain>",
	Aliases: []string{"defaults"},
	Short:   "Default gets the blockchain params.",
	Long: `
Get the blockchain specific parameters for a deployed blockchain.

Response: The params as a list of key value params, of name and type respectively
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		util.JsonRpcCallAndPrint("get_defaults", args)
	},
}

var getConfigsCmd = &cobra.Command{
	Use:     "configs <blockchain> [file]",
	Aliases: []string{"config"},
	Short:   "Get the resources for a blockchain",
	Long: `
Get the resources for a blockchain. With one argument, lists what is available. With two
	arguments, get the contents of the file

Params: The blockchain to get the resources of, the resource/file name 

Response: The resoures as a list of key value params, of name and type respectively
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		util.JsonRpcCallAndPrint("get_resources", args)
	},
}

var getStatsCmd = &cobra.Command{
	Use:   "stats <command>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics
	`,
	Run: util.PartialCommand,
}

var statsByTimeCmd = &cobra.Command{
	Use:   "time <start time> <end time>",
	Short: "Get stastics by time",
	Long: `
Stats time will allow the user to get statistics by specifying a start time and stop time (unix time stamp).

Params: start unix timestamp, end unix timestamp

Response: JSON representation of network statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		util.JsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  util.CheckAndConvertInt64(args[0], "start unix timestamp"),
			"endTime":    util.CheckAndConvertInt64(args[1], "end unix timestamp"),
			"startBlock": 0,
			"endBlock":   0,
		})
	},
}

var statsByBlockCmd = &cobra.Command{
	Use:   "block <start block> <end block>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats block will allow the user to get statistics regarding the network.

Response: JSON representation of statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		util.JsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  0,
			"endTime":    0,
			"startBlock": util.CheckAndConvertInt64(args[0], "start block number"),
			"endBlock":   util.CheckAndConvertInt64(args[1], "end block number"),
		})
	},
}

var statsPastBlocksCmd = &cobra.Command{
	Use:   "past <number of blocks> ",
	Short: "Get stastics of a blockchain from the past x blocks",
	Long: `
Stats block will allow the user to get statistics regarding the network.

Response: JSON representation of statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		util.JsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  0,
			"endTime":    0,
			"startBlock": util.CheckAndConvertInt64(args[0], "blocks") * -1, //Negative number signals past
			"endBlock":   0,
		})
	},
}

var statsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all stastics of a blockchain",
	Long: `
Stats all will allow the user to get all the statistics regarding the network.

Response: JSON representation of network statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("all_stats", []string{})
	},
}

/*
Work underway on generalized use commands to consolidate all the different
commands separated by blockchains.
*/

func getBlockCobra(cmd *cobra.Command, args []string) {
	util.CheckArguments(cmd, args, 1, 1)
	blockNum := util.CheckAndConvertInt(args[0], "block number")

	if blockNum < 1 {
		util.PrintStringError("Unable to get block information from block 0. Please provide a block number greater than 0.")
		os.Exit(1)
	}
	/*res, err := util.JsonRpcCall("get_block_number", []string{})
	if err != nil {
		util.PrintErrorFatal(err)
	}

	blocknum := int(res.(float64))
	if blocknum < 1 {
		util.PrintStringError("Unable to get block information because no blocks have been created." +
			" Please use the command 'whiteblock miner start' to start generating blocks.")
		os.Exit(1)
	}*/

	res, err := util.JsonRpcCall("get_block", args)
	if err != nil { //try a few nodes
		nodes, er := GetNodes()
		if er != nil {
			util.PrintErrorFatal(er)
		}
		for i := range nodes {
			res, err = util.JsonRpcCall("get_block", []interface{}{args[0], i})
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		util.PrintErrorFatal(err)
	}
	cmd.Println(util.Prettypi(res))
}

var getBlockCmd = &cobra.Command{
	Use:   "block <command>",
	Short: "Get information regarding blocks",
	Run:   getBlockCobra,
}

var getBlockNumCmd = &cobra.Command{
	Use:   "number",
	Short: "Get the block number",
	Long: `
Gets the most recent block number that had been added to the blockchain.

Response: block number
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("get_block_number", []string{})
	},
}

var getBlockInfoCmd = &cobra.Command{
	Use:   "info <block number>",
	Short: "Get the information of a block",
	Long: `
Gets the information inside a block including transactions and other information relevant to the currently connected blockchain.

Response: JSON representation of the block
	`,
	Run: getBlockCobra,
}

var getTxCmd = &cobra.Command{
	Use:   "tx <command>",
	Short: "Get information regarding transactions",
	Run:   util.PartialCommand,
}

var getTxRecentCmd = &cobra.Command{
	Use:   "recent [number of tx]",
	Short: "Get transaction information",
	Long: `Get the tx hash(es) of recently sent transactions
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		var num int = 5
		if len(args) > 0 {
			num = util.CheckAndConvertInt(args[0], "number of transactions")
		}

		util.JsonRpcCallAndPrint("state::get_recent_tx", []interface{}{num})
	},
}

var getTxInfoCmd = &cobra.Command{
	Use:   "info <tx hash>",
	Short: "Get transaction information",
	Long: `
Get a transaction by its hash. The user can find the transaction hash by viewing block information. 
To view block information, the command 'get block info <block number>' can be used.

Response: JSON representation of the transaction.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		util.JsonRpcCallAndPrint("get_transaction", args)
	},
}

var getTxReceiptCmd = &cobra.Command{
	Use:   "receipt <tx hash>",
	Short: "Get the transaction receipt",
	Long: `
Get the transaction receipt by the tx hash. The user can find the transaction hash by viewing block information. 
To view block information, the command 'get block info <block number>' can be used.

Response: JSON representation of the transaction receipt.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		util.JsonRpcCallAndPrint("get_transaction_receipt", args)
	},
}

var getAccountCmd = &cobra.Command{
	Aliases: []string{"accounts"},
	Use:     "account",
	Short:   "Get account information",
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("state::get", []string{"accounts"})
	},
}

var getAccountInfoCmd = &cobra.Command{
	// Hidden: true,
	Use:   "info",
	Short: "Get account information",
	Long: `
Gets the account information relevant to the currently connected blockchain.

Response: JSON representation of the accounts information.
`,
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("accounts_status", []string{})
	},
}

var getContractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Get contracts deployed to network.",
	Long: `
Gets the list of contracts that were deployed to the network. The information includes the address that deployed the contract, 
the contract name, and the contract's address.

Response: JSON representation of the contract information.
`,
	Run: func(cmd *cobra.Command, args []string) {

		contracts, err := readContractsFile()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if len(contracts) == 0 {
			util.PrintStringError("No smart contract has been deployed yet." +
				" Please use the command 'whiteblock geth solc deploy <smart contract> to deploy a smart contract.")
			os.Exit(1)
		} else {
			fmt.Println(util.Prettyp(string(contracts)))
		}
	},
}

func init() {
	getNodesCmd.Flags().Bool("all", false, "output all of the nodes, even if they are no longer running")
	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd, getDefaultsCmd,
		getSupportedCmd, getRunningCmd, getConfigsCmd, getTestnetIDCmd, getBuildCmd)

	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsPastBlocksCmd, statsAllCmd)

	getCmd.AddCommand(getBlockCmd, getTxCmd, getAccountCmd, getContractsCmd)
	getBlockCmd.AddCommand(getBlockNumCmd, getBlockInfoCmd)
	getTxCmd.AddCommand(getTxInfoCmd, getTxReceiptCmd, getTxRecentCmd)
	getAccountCmd.AddCommand(getAccountInfoCmd)

	RootCmd.AddCommand(getCmd)
}
