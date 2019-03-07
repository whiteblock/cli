package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	util "../util"
)

/*
func cwFile(path, data string) {
	time := time.Now().UTC().String()
	time = strings.Replace(time, " ", "", -1)

	file, err := os.Create(path + "/dataset_" + time + ".txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}
*/

func GetNodes() ([]Node, error) {
	res, err := jsonRpcCall("nodes", []string{})
	if err != nil {
		return nil, err
	}
	tmp := res.([]interface{})
	nodes := []map[string]interface{}{}
	for _, t := range tmp {
		nodes = append(nodes, t.(map[string]interface{}))
	}

	out := []Node{}
	for _, node := range nodes {
		out = append(out, Node{
			LocalID:   int(node["localId"].(float64)),
			Server:    int(node["server"].(float64)),
			TestNetID: node["testNetId"].(string),
			ID:        node["id"].(string),
			IP:        node["ip"].(string),
			Label:     node["label"].(string),
		})
	}
	return out, nil
}

var logTail int = -1

func readContractsFile() ([]byte, error) {
	cwd := os.Getenv("HOME")
	b, err := ioutil.ReadFile(cwd + "/smart-contracts/whiteblock/contracts.json")
	if err != nil {
		//fmt.Print(err)
	}
	return b, nil
}

var getCmd = &cobra.Command{
	Use:   "get <command>",
	Short: "Get server and network information.",
	Long:  "\nGet will ouput server and network information and statstics.\n",
	Run:   util.PartialCommand,
}

var getServerCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"servers"},
	Short:   "Get server information.",
	Long:    "\nServer will ouput server information.\n",
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("get_servers", []string{})
	},
}

var getNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Nodes will show all nodes in the network.",
	Long: `
Nodes will output all of the nodes in the current network.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("get_nodes", []string{})
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
		util.CheckArguments(args, 0, 0)
		jsonRpcCallAndPrint("state::is_running", []string{})
		jsonRpcCallAndPrint("state::what_is_running", []string{})
	},
}

var getLogCmd = &cobra.Command{
	Use:   "log <node>",
	Short: "Log will dump data pertaining to the node.",
	Long: `
Get stdout and stderr from a node.

Params: node number

Response: stdout and stderr of the blockchain process
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		previousBuild,err := getPreviousBuild()
		if err != nil{
			util.PrintErrorFatal(err)
		}
		s := previousBuild.Servers[0]
		n, err := strconv.Atoi(args[0])

		if err != nil {
			util.InvalidInteger("node", args[0], true)
		}

		jsonRpcCallAndPrint("log", map[string]int{
			"server": s,
			"node":   n,
			"lines":  logTail,
		})

	},
}

var getDefaultsCmd = &cobra.Command{
	Use:   "default <blockchain>",
	Short: "Default gets the blockchain params.",
	Long: `
Get the blockchain specific parameters for a deployed blockchain.

Params: <blockchain>
Format: The blockchain to get the build params of

Response: The params as a list of key value params, of name and type respectively
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("get_defaults", args)
	},
}


var getConfigsCmd = &cobra.Command{
	Use:   "configs <blockchain>",
	Short: "Get the resources for a blockchain",
	Long: `
Get the resources for a blockchain

Params: <blockchain>
Format: The blockchain to get the resources of

Response: The resoures as a list of key value params, of name and type respectively
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("get_resources", args)
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

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

Response: JSON representation of network statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 2, 2)
		jsonRpcCallAndPrint("stats", map[string]int64{
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

Params: Block numbers
Format: <start block number> <end block number>

Response: JSON representation of statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 2, 2)
		jsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  0,
			"endTime":    0,
			"startBlock": util.CheckAndConvertInt64(args[0], "start block number"),
			"endBlock":   util.CheckAndConvertInt64(args[1], "end block number"),
		})
	},
}

var statsPastBlocksCmd = &cobra.Command{
	Use:   "past <blocks> ",
	Short: "Get stastics of a blockchain from the past x blocks",
	Long: `
Stats block will allow the user to get statistics regarding the network.

Params: Number of blocks 
Format: <blocks>

Response: JSON representation of statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("stats", map[string]int64{
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
		jsonRpcCallAndPrint("all_stats", []string{})
	},
}

/*
Work underway on generalized use commands to consolidate all the different
commands separated by blockchains.
*/

var getBlockCmd = &cobra.Command{
	// Hidden: true,
	Use:   "block <command>",
	Short: "Get information regarding blocks",
	Run:   util.PartialCommand,
}

var getBlockNumCmd = &cobra.Command{
	// Hidden: true,
	Use:   "number",
	Short: "Get the block number",
	Long: `
Gets the most recent block number that had been added to the blockchain.

Response: block number
	`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("get_block_number", []string{})
	},
}

var getBlockInfoCmd = &cobra.Command{
	// Hidden: true,
	Use:   "info <block number>",
	Short: "Get the information of a block",
	Long: `
Gets the information inside a block including transactions and other information relevant to the currently connected blockchain.

Format: <Block Number>
Params: Block number

Response: JSON representation of the block
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		blockNum := 0
		var err error
		if len(args) > 0 {
			blockNum, err = strconv.Atoi(args[0])
			if err != nil {
				util.PrintStringError("Invalid block number formatting.")
				return
			}
		}
		if blockNum < 1 && len(args) > 0 {
			util.PrintStringError("Unable to get block information from block 0. Please provide a block number greater than 0.")
			return
		} else {
			res, err := jsonRpcCall("get_block_number", []string{})
			if err != nil {
				util.PrintErrorFatal(err)
			}
			blocknum := int(res.(float64))
			if blocknum < 1 {
				util.PrintStringError("Unable to get block information because no blocks have been created. Please use the command 'whiteblock miner start' to start generating blocks.")
				return
			}
		}
		jsonRpcCallAndPrint("get_block", args)
	},
}

var getTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "tx <command",
	Short: "Get information regarding transactions",
	Run:   util.PartialCommand,
}

var getTxInfoCmd = &cobra.Command{
	// Hidden: true,
	Use:   "info <tx hash>",
	Short: "Get transaction information",
	Long: `
Get a transaction by its hash. The user can find the transaction hash by viewing block information. To view block information, the command 'get block info <block number>' can be used.

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("get_transaction", args)
	},
}

// eth::get_transaction_receipt does not work.
/*
var getTxReceiptCmd = &cobra.Command{
	// Hidden: true,
	Use:   "receipt <tx hash>",
	Short: "Get the transaction receipt",
	Long: `
Get the transaction receipt by the tx hash. The user can find the transaction hash by viewing block information. To view block information, the command 'get block info <block number>' can be used.

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction receipt.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		command := ""
		switch blockchain {
		case "ethereum":
			command = "eth::get_transaction_receipt"
		case "eos":
			fmt.Println("This function is not supported for the syscoin client.")
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		jsonRpcCallAndPrint(command, args)
	},
}
*/

var getAccountCmd = &cobra.Command{
	// Hidden: true,
	Use:   "account <command>",
	Short: "Get account information",
	Run:   util.PartialCommand,
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
		jsonRpcCallAndPrint("accounts_status", []string{})
	},
}

var getContractsCmd = &cobra.Command{
	// Hidden: true,
	Use:   "contracts",
	Short: "Get contracts deployed to network.",
	Long: `
Gets the list of contracts that were deployed to the network. The information includes the address that deployed the contract, the contract name, and the contract's address.

Response: JSON representation of the contract information.
`,
	Run: func(cmd *cobra.Command, args []string) {

		contracts, err := readContractsFile()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(contracts) == 0 {
			util.PrintStringError("No smart contract has been deployed yet. Please use the command 'whiteblock geth solc deploy <smart contract> to deploy a smart contract.")
		} else {
			fmt.Println(prettyp(string(contracts)))
		}
		
		
	},
}

func init() {
	getCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	getServerCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	getNodesCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getStatsCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsByTimeCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsByBlockCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsAllCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getLogCmd.Flags().IntVarP(&logTail, "tail", "t", -1, "see tail(1)")

	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd, getDefaultsCmd, getRunningCmd, getLogCmd,getConfigsCmd)
	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsPastBlocksCmd, statsAllCmd)

	// dev commands that are currently being implemented
	getCmd.AddCommand(getBlockCmd, getTxCmd, getAccountCmd, getContractsCmd)
	getBlockCmd.AddCommand(getBlockNumCmd, getBlockInfoCmd)
	getTxCmd.AddCommand(getTxInfoCmd)
	getAccountCmd.AddCommand(getAccountInfoCmd)

	RootCmd.AddCommand(getCmd)
}
