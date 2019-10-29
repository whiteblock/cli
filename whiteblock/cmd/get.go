package cmd

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
)

func GetNodes() []Node {
	var out []Node
	err := util.JsonRpcCallP("nodes", []string{build.GetPreviousBuildID()}, &out)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	log.WithFields(log.Fields{"nodes": out}).Trace("raw nodes")
	return out
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
		util.Print(build.GetPreviousBuildID())
	},
}

var getBuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"built"},
	Short:   "Get the last applied build",
	Long:    "\nGet the last applied build.\n",
	Run: func(cmd *cobra.Command, args []string) {
		prevBuild, err := build.GetPreviousBuild()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		util.Print(prevBuild)
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
		util.Print([]string(sortedBlockchains))
	},
}

var getNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Nodes will show all nodes in the network.",
	Long:    "\nNodes will output all of the nodes in the current network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		testnetID := build.GetPreviousBuildID()
		if util.GetBoolFlagValue(cmd, "all") {
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
		util.Print(out)
	},
}

var getNodesExternalCmd = &cobra.Command{
	Use:     "external",
	Aliases: []string{"externals", "exposed"},
	Short:   "Get the port mappings of the nodes",
	Long:    "\nGet the port mappings of the nodes\n",

	Run: func(cmd *cobra.Command, args []string) {
		nodes := GetNodes()
		out := []map[string]interface{}{}
		for _, node := range nodes {
			out = append(out, map[string]interface{}{
				"id":           node.ID,
				"absNum":       node.AbsoluteNum,
				"portMappings": node.PortMappings,
				"protocol":     node.Protocol,
			})
		}
		util.Print(out)
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
	Get some general statistics on the blockchain network.
	Note: This is will be behind by a few blocks to ensure accuracy.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("all_stats", []string{})
	},
}

/*
Work underway on generalized use commands to consolidate all the different
commands separated by blockchains.
*/
func getBlockJsonRpcCall(rpc string, args []string) {
	res, err := util.JsonRpcCall(rpc, args)
	if err != nil { //try a few nodes
		nodes := GetNodes()
		for i := range nodes {
			res, err = util.JsonRpcCall(rpc, []interface{}{args[0], i})
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		util.PrintErrorFatal(err)
	}
	util.Print(res)
}

func getBlockCobra(cmd *cobra.Command, args []string) {
	util.CheckArguments(cmd, args, 1, 1)
	out, err := strconv.ParseInt(args[0], 0, 32) //Check if the input is an integer-> block number
	if err != nil {
		// Block Hash
		getBlockJsonRpcCall("get_block_by_hash", args)
	} else {
		// Block number
		blockNum := int(out)
		if blockNum < 1 {
			util.PrintErrorFatal("Unable to get block information from block 0. Please provide a block number greater than 0.")
		}
		getBlockJsonRpcCall("get_block", args)
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
}

func getBlockHeightByNode(cmd *cobra.Command, args []string) {
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	util.CheckArguments(cmd, args, 0, 1)

	if util.GetBoolFlagValue(cmd, "all") {
		nodes := len(GetNodes())

		blockHeights := make([]string, nodes)

		for i := 0; i < nodes; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				res, err := util.JsonRpcCall("get_block_number", []interface{}{i})
				if err != nil {
					util.PrintErrorFatal(err)
				}

				mux.Lock()
				blockHeights[i] = fmt.Sprintf("Node %v: %v", i, res)
				mux.Unlock()
			}(i)
		}
		wg.Wait()

		util.Print(blockHeights)

		return
	}

	if len(args) == 0 {
		util.JsonRpcCallAndPrint("get_block_number", []interface{}{0})
		return
	}

	util.JsonRpcCallAndPrint("get_block_number", []interface{}{args[0]})

	return
}

func getPrivateKeys(cmd *cobra.Command, args []string) {
	res, err := util.JsonRpcCall("state::info", args)
	if err != nil {
		util.PrintErrorFatal(err)
		return
	}

	privKeys := make([]string, 0)

	for _, keys := range res.(map[string]interface{}) {
		if reflect.DeepEqual(reflect.TypeOf(keys).String(), "map[string]interface {}") { // TODO is there a cleaner way to do this?
			for i, val := range keys.(map[string]interface{}) {
				if i == "privateKey" {
					privKeys = append(privKeys, val.(string))
				}
			}
		}
	}

	util.Print(privKeys)
	return
}

var getBlockCmd = &cobra.Command{
	Use:   "block <command>",
	Short: "Get information regarding blocks",
	Run:   getBlockCobra,
}

var getBlockNumCmd = &cobra.Command{
	Use:   "number [node]",
	Short: "Get the block number of a single node or use --all to get block heights of all nodes",
	Long: `
Gets the most recent block number that had been added to the blockchain.

Response: block number
	`,
	Run: getBlockHeightByNode,
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
		out := []interface{}{}
		stepSize := 50
		i := 0
		for {
			res := []interface{}{}
			err := util.JsonRpcCallP("state::get_page", []interface{}{"accounts", i, stepSize}, &res)
			if err != nil {
				util.PrintErrorFatal(err)
			}
			out = append(out, res...)
			log.WithFields(log.Fields{"i": i, "stepSize": stepSize, "results": len(res), "total": len(out)}).Debug("got some accounts")
			if len(res) != stepSize {

				break
			}
			i += stepSize
		}
		util.Print(out)
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

var getPrivateKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Get private keys",
	Long: `
Gets the private keys of _______________________________.
`,
	Run: getPrivateKeys,
}

var getBiomeCmd = &cobra.Command{
	Use:   "biome",
	Short: "Get the biome id",
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("get_biome_id", []string{})
	},
}

var getBoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Get box information",
	Run: func(cmd *cobra.Command, args []string) {
		util.Print(conf.ServerAddr)
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
		var contracts []interface{}
		err := util.ReadTestnetStore("contracts", &contracts)
		if err != nil {
			util.PrintErrorFatal("No smart contract has been deployed yet." +
				" Please use the command 'whiteblock geth solc deploy <smart contract> to deploy a smart contract.")

		}
		util.Print(contracts)
	},
}

func init() {
	getNodesCmd.AddCommand(getNodesExternalCmd)
	getNodesCmd.Flags().Bool("all", false, "output all of the nodes, even if they are no longer running")

	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd, getDefaultsCmd,
		getSupportedCmd, getRunningCmd, getConfigsCmd, getTestnetIDCmd, getBuildCmd, getPrivateKeysCmd, getBoxCmd)

	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsPastBlocksCmd, statsAllCmd)

	getCmd.AddCommand(getBlockCmd, getTxCmd, getAccountCmd, getContractsCmd, getBiomeCmd)

	getBlockCmd.AddCommand(getBlockNumCmd, getBlockInfoCmd)
	getBlockNumCmd.Flags().BoolP("all", "a", false, "output block heights of all nodes")

	getTxCmd.AddCommand(getTxInfoCmd, getTxReceiptCmd, getTxRecentCmd)

	getAccountCmd.AddCommand(getAccountInfoCmd)

	RootCmd.AddCommand(getCmd)
}
