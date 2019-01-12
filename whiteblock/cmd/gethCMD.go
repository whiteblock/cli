package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var (
	gethcommand string
)

var gethCmd = &cobra.Command{
	Use:   "geth <command>",
	Short: "Run geth commands",
	Long: `
Geth will allow the user to get infromation and run geth commands.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		return
	},
}

var gethConsole = &cobra.Command{
	Use:   "console <node number>",
	Short: "Logs into the geth console",
	Long: `
Console will log into the geth console.

Response: stdout of geth console`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "nodes"
		out1 := []byte(wsEmitListen(serverAddr, command1, ""))
		var node Node
		json.Unmarshal(out1, &node)
		nodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}

		command2 := "exec"
		param := "{\"server\":" + server + ",\"node\":" + args[0] + ",\"command\":\"service ssh start\"}"
		wsEmitListen(serverAddr, command2, param)

		err = unix.Exec("/usr/bin/ssh", []string{"ssh", "-o", "StrictHostKeyChecking no", "root@" + fmt.Sprintf(node[nodeNumber].IP), "tmux", "attach", "-t", "whiteblock"}, os.Environ())
		log.Fatal(err)
	},
}

var getBlockNumberCmd = &cobra.Command{
	Use:   "get_block_number",
	Short: "Get block number",
	Long: `
Get the current highest block number of the chain

Response: The block number`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_block_number"
		param := ""
		// fmt.Println(command, param)
		if len(args) >= 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getBlockCmd = &cobra.Command{
	Use:   "get_block <block number>",
	Short: "Get block information",
	Long: `
Get the data of a block

Format: <Block Number>
Params: Block number

Response: JSON Representation of the block.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_block"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) < 1 || len(args) > 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getAccountCmd = &cobra.Command{
	Use:   "get_accounts",
	Short: "Get account information",
	Long: `
Get a list of all unlocked accounts

Response: A JSON array of the accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_accounts"
		param := ""
		// fmt.Println(command)
		if len(args) >= 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getBalanceCmd = &cobra.Command{
	Use:   "get_balance <address>",
	Short: "Get account balance information",
	Long: `
Get the current balance of an account

Format: <address>
Params: Account address

Response: The integer balance of the account in wei`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_balance"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) < 1 || len(args) > 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var sendTxCmd = &cobra.Command{
	Use:   "send_transaction <from address> <to address> <gas> <gas price> <value to send>",
	Short: "Sends a transaction",
	Long: `
Send a transaction between two accounts

Format: <from> <to> <gas> <gas price> <value>
Params: Sending account, receiving account, gas, gas price, amount to send in ETH

Response: The transaction hash`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(command)
		if len(args) != 5 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::send_transaction"
		param := strings.Join(args[:], " ")

		weiToInt, err := strconv.Atoi(args[4])
		if err != nil {
			panic(err)
		}
		weiToEth := weiToInt * 1000000000000000000
		args[1] = strconv.Itoa(weiToEth)

		wsEmitListen(serverAddr, command, param)
	},
}

var getTxCountCmd = &cobra.Command{
	Use:   "get_transaction_count <address> [block number]",
	Short: "Get transaction count",
	Long: `
Get the transaction count sent from an address, optionally by block

Format: <address> [block number]
Params: The sender account, a block number

Response: The transaction count`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_transaction_count"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) < 1 || len(args) > 2 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getTxCmd = &cobra.Command{
	Use:   "get_transaction <hash>",
	Short: "Get transaction information",
	Long: `
Get a transaction by its hash

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_transaction"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) < 1 || len(args) > 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getTxReceiptCmd = &cobra.Command{
	Use:   "get_transaction_receipt <hash>",
	Short: "Get transaction receipt",
	Long: `
Get the transaction receipt by the tx hash

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction receipt.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_transaction_receipt"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) < 1 || len(args) > 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getHashRateCmd = &cobra.Command{
	Use:   "get_hash_rate",
	Short: "Get hash rate",
	Long: `
Get the current hash rate per node

Response: The hash rate of a single node in the network`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_hash_rate"
		param := ""
		// fmt.Println(command)
		if len(args) >= 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var startTxCmd = &cobra.Command{
	Use:   "start_transactions <tx/s> <value> [destination]",
	Short: "Start transactions",
	Long: `
Start sending transactions according to the given parameters, value = -1 means randomize value.

Format: <tx/s> <value> [destination]
Params: The amount of transactions to send in a second, the value of each transaction in wei, the destination for the transaction
`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::start_transactions"
		// fmt.Println(command)
		if len(args) <= 1 || len(args) > 3 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		weiToInt, err := strconv.Atoi(args[1])
		weiToEth := weiToInt * 1000000000000000000
		args[1] = strconv.Itoa(weiToEth)
		if err != nil {
			panic(err)
		}
		param := strings.Join(args[:], " ")
		wsEmitListen(serverAddr, command, param)
	},
}

var stopTxCmd = &cobra.Command{
	Use:   "stop_transactions",
	Short: "Stop transactions",
	Long: `
Stops the sending of transactions if transactions are currently being sent`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::stop_transactions"
		param := ""
		// fmt.Println(command)
		if len(args) >= 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var startMiningCmd = &cobra.Command{
	Use:   "start_mining [node 1 number] [node 2 number]...",
	Short: "Start Mining",
	Long: `
Send the start mining signal to nodes, may take a while to take effect due to DAG generation

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to start mining or None for all nodes

Response: The number of nodes which successfully received the signal to start mining`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::start_mining"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		wsEmitListen(serverAddr, command, param)
	},
}

var stopMiningCmd = &cobra.Command{
	Use:   "stop_mining [node 1 number] [node 2 number]...",
	Short: "Stop mining",
	Long: `
Send the stop mining signal to nodes

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to stop mining or None for all nodes


Response: The number of nodes which successfully received the signal to stop mining`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::stop_mining"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		wsEmitListen(serverAddr, command, param)
	},
}

var blockListenerCmd = &cobra.Command{
	Use:   "block_listener [block number]",
	Short: "Get block listener",
	Long: `
Get all blocks and continue to subscribe to new blocks

Format: [block number]
Params: The block number to start at or None for all blocks

Response: Will emit on eth::block_listener for every block after the given block or 0 that exists/has been created`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::block_listener"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) > 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

var getRecentSentTxCmd = &cobra.Command{
	Use:   "get_recent_sent_tx [number]",
	Short: "Get recently sent transaction",
	Long: `
Get a number of the most recent transactions sent

Format: [number]
Params: The number of transactions to retrieve

Response: JSON object of transaction data`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "eth::get_recent_sent_tx"
		param := strings.Join(args[:], " ")
		// fmt.Println(command)
		if len(args) > 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		wsEmitListen(serverAddr, command, param)
	},
}

func init() {
	// gethCmd.Flags().StringVarP(&gethcommand, "command", "c", "", "Geth command")
	gethCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	//geth subcommands
	gethCmd.AddCommand(getBlockNumberCmd, getBlockCmd, getAccountCmd, getBalanceCmd, sendTxCmd, getTxCountCmd, getTxCmd, getTxReceiptCmd,
		getHashRateCmd, startTxCmd, stopTxCmd, startMiningCmd, stopMiningCmd, blockListenerCmd, getRecentSentTxCmd)

	RootCmd.AddCommand(gethCmd)
}
