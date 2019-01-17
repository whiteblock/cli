package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

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

var getCmd = &cobra.Command{
	Use:   "get <command>",
	Short: "Get server and network information.",
	Long: `
Get will ouput server and network information and statstics.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.\n")
		cmd.Help()
		return
	},
}

var getServerCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"servers"},
	Short:   "Get server information.",
	Long: `
Server will ouput server information.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_servers"

		fmt.Println(prettyp(wsEmitListen(serverAddr, command, "")))
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_nodes"
		fmt.Println(wsEmitListen(serverAddr, command, ""))
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
		if len(args) != 0 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "state::is_running"
		command2 := "state::what_is_running"
		fmt.Println(wsEmitListen(serverAddr, command1, ""))
		fmt.Println(wsEmitListen(serverAddr, command2, ""))

	},
}

var getLogCmd = &cobra.Command{
	Use:   "log <node number>",
	Short: "Log will dump data pertaining to the node.",
	Long: `
Get stdout and stderr from a node.

Params: node number

Response: stdout and stderr of the blockchain process
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "log"
		param := "{\"server\":" + fmt.Sprintf(server) + ",\"node\":" + args[0] + "}"
		fmt.Println(wsEmitListen(serverAddr, command, param))
	},
}

var getNetworkDefaultsCmd = &cobra.Command{
	Use:   "default <blockchain>",
	Short: "Default gets the blockchain params.",
	Long: `
Get the blockchain specific parameters for a deployed blockchain.

Params: <blockchain>
Format: The blockchain to get the build params of

Response: The params as a list of key value params, of name and type respectively
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_defaults"
		fmt.Println(wsEmitListen(serverAddr, command, args[0]))
	},
}

var getStatsCmd = &cobra.Command{
	Use:   "stats <command>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nError: Invalid number of arguments given\n")
		cmd.Help()
		return
	},
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
		if len(args) != 2 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":" + args[0] + ",\"endTime\":" + args[1] + ",\"startBlock\":0,\"endBlock\":0}"
		data := wsEmitListen(serverAddr, command, param)
		fmt.Println(data)
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
		if len(args) != 2 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":0,\"endTime\":0,\"startBlock\":" + args[0] + ",\"endBlock\":" + args[1] + "}"
		data := wsEmitListen(serverAddr, command, param)
		fmt.Println(data)
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
		if len(args) != 1 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"

		blocks, err := strconv.Atoi(args[0])
		if err != nil {
			InvalidArgument(args[0])
			cmd.Help()
			return
		}
		blocks *= -1
		param := fmt.Sprintf("{\"startTime\":0,\"endTime\":0,\"startBlock\":%d,\"endBlock\":0}", blocks)
		data := wsEmitListen(serverAddr, command, param)
		fmt.Println(data)
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "all_stats"
		data := wsEmitListen(serverAddr, command, "")
		fmt.Println(data)
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.")
		cmd.Help()
		return
	},
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := ""
		switch blockchain {
		case "ethereum":
			command = "eth::get_block_number"
		case "eos":
			command = "eos::get_block_number"
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		data := wsEmitListen(serverAddr, command, "")
		fmt.Println(data)
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := ""
		switch blockchain {
		case "ethereum":
			command = "eth::get_block"
		case "eos":
			command = "eos::get_block"
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		data := wsEmitListen(serverAddr, command, args[0])
		fmt.Println(data)
	},
}

var getTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "tx <command",
	Short: "Get information regarding transactions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.")
		cmd.Help()
		return
	},
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := ""
		switch blockchain {
		case "ethereum":
			command = "eth::get_transaction"
		case "eos":
			fmt.Println("This function is not supported for the eos client.")
			return
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		data := wsEmitListen(serverAddr, command, args[0])
		fmt.Println(data)
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
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
		data := wsEmitListen(serverAddr, command, args[0])
		fmt.Println(data)
	},
}
*/

var getAccountCmd = &cobra.Command{
	// Hidden: true,
	Use:   "account <command>",
	Short: "Get account information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nNo command given. Please choose a command from the list below.")
		cmd.Help()
		return
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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := ""
		switch blockchain {
		case "ethereum":
			command = "eth::accounts_status"
			data := wsEmitListen(serverAddr, command, "")
			fmt.Println(data)
		case "eos":
			nodenum, _ := strconv.Atoi(nodes)
			AccBalances := make([]interface{}, 0)
			for i := 0; i < nodenum; i++ {
				AccBalances = append(AccBalances, wsEmitListen(serverAddr, "eos::get_info", strconv.Itoa(i)))
			}
			fmt.Println(AccBalances)
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
			return
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
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

	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd, getNetworkDefaultsCmd, getRunningCmd, getLogCmd)
	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsPastBlocksCmd, statsAllCmd)

	// dev commands that are currently being implemented
	getCmd.AddCommand(getBlockCmd, getTxCmd, getAccountCmd)
	getBlockCmd.AddCommand(getBlockNumCmd, getBlockInfoCmd)
	getTxCmd.AddCommand(getTxInfoCmd)
	getAccountCmd.AddCommand(getAccountInfoCmd)

	RootCmd.AddCommand(getCmd)
}
