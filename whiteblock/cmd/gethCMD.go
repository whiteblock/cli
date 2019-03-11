package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"os/exec"
	"strconv"
	"strings"

	util "../util"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var (
	gethcommand string
)

type Contracts struct {
	DeployedNodeAddress string `json:",omitempty"`
	ContractName        string `json:,omitempty`
	ContractAddress     string `json:",omitempty"`
}

const (
	exSolFile = `pragma solidity ^0.4.25;
	contract helloWorld {
		function renderHelloWorld () public pure returns (string memory) {
			return "helloWorld";
		}
	}
	`
	deployJs = `//deploy.js
const Web3 = require('web3');
const { abi, bytecode } = require('./compile');
const web3 = new Web3(new Web3.providers.HttpProvider("http://"+process.argv[3].toString()+":8545"));
const deploy = async () => {
	const accounts = await web3.eth.getAccounts();
	console.log('Attempting to deploy from account', accounts[0]);
	const result = await new web3.eth.Contract(abi).deploy({ data: '0x' + bytecode}).send({ gas: '1000000', from: accounts[0] });
	console.log('Contract deployed to', result.options.address);
};
deploy();`

	compileJS = `//compile.js
const path = require('path');
const fs = require('fs');
const solc = require('solc');
const contractName = process.argv[2].split('.')[0];
const helloWorldPath = path.resolve(__dirname, process.argv[2]);
const input = fs.readFileSync(helloWorldPath);
const output = solc.compile(input.toString().toLowerCase(), 1);
const bytecode = output.contracts[':'+contractName].bytecode;
const abi = JSON.parse(output.contracts[':'+contractName].interface);
module.exports = {abi, bytecode};`
)

func writeContractListFile(addrs string) {
	cwd := os.Getenv("HOME")
	err := os.MkdirAll(cwd+"/smart-contracts/whiteblock/", 0755)
	if err != nil {
		log.Fatalf("could not create directory: %s", err)
	}

	file, err := os.OpenFile(cwd+"/smart-contracts/whiteblock/contracts.json", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	newCont := ""
	cont, err := ioutil.ReadFile(cwd + "/smart-contracts/whiteblock/contracts.json")
	if err != nil {
		fmt.Println(err)
	}
	contents := string(cont)
	if len(cont) > 0 {
		contents = strings.TrimRight(contents, "]")
		addrs = strings.TrimLeft(addrs, "[")
		addrs = ", " + addrs
		newCont = contents + addrs
		_, err = file.WriteString(newCont)
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
	} else {
		_, err = file.WriteString(addrs)
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
	}

}

func checkContractDir() {
	cwd := os.Getenv("HOME")
	if _, err := os.Stat(cwd + "/smart-contracts/"); os.IsNotExist(err) {
		fmt.Println("'smart-contracts' directory could not be found. Creating the directory 'smart-contracts' in home directory.")
		fmt.Println("Preparing the dependencies to deploy smart contracts.")
		err := os.MkdirAll(cwd+"/smart-contracts/", 0755)
		if err != nil {
			log.Fatalf("could not create directory: %s", err)
			return
		}
		solFile, err := os.Create(cwd + "/smart-contracts/helloworld.sol")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
			return
		}
		solFile.Write([]byte(exSolFile))
		compileFile, err := os.Create(cwd + "/smart-contracts/compile.js")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
			return
		}
		compileFile.Write([]byte(compileJS))
		deployFile, err := os.Create(cwd + "/smart-contracts/deploy.js")
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		deployFile.Write([]byte(deployJs))

		defer solFile.Close()
		defer compileFile.Close()
		defer deployFile.Close()

		return
	}
}

func checkContractFiles(fileName string) bool {
	cwd := os.Getenv("HOME")
	if _, err := os.Stat(cwd + "/smart-contracts/node_modules"); err != nil {
		fmt.Println("Smartcontracts have not been initialized. Please run 'geth solc init' to deploy a smart contract.")
		return false
	}
	if _, err := os.Stat(cwd + "/smart-contracts/package.json"); err != nil {
		fmt.Println("Smartcontracts have not been initialized. Please run 'geth solc init' to deploy a smart contract.")
		return false
	}
	if _, err := os.Stat(cwd + "/smart-contracts/compile.js"); err != nil {
		fmt.Println("Smartcontracts have not been initialized. Please run 'geth solc init' to deploy a smart contract.")
		return false
	}
	if _, err := os.Stat(cwd + "/smart-contracts/deploy.js"); err != nil {
		fmt.Println("Smartcontracts have not been initialized. Please run 'geth solc init' to deploy a smart contract.")
		return false
	}
	if _, err := os.Stat(cwd + "/smart-contracts/" + fileName); err != nil {
		fmt.Println(fileName + " is not in the directory 'smart-contracts'. Please make sure that the file is located in the directory.")
		return false
	}
	return true
}

func installNpmDeps() {
	cwd := os.Getenv("HOME")
	if _, err := os.Stat(cwd + "/smart-contracts/node_modules"); err == nil {
		return
	}

	fmt.Printf("\rDependencies are being loaded...")

	npmInitCmd := exec.Command("npm", "init", "-y")
	npmInitCmd.Dir = cwd + "/smart-contracts/"
	_, err := npmInitCmd.Output()
	if err != nil {
		log.Println(err)
	}
	// fmt.Printf("%s", output)

	npmInstWeb3Cmd := exec.Command("npm", "install", "web3@1.0.0-beta.31")
	npmInstWeb3Cmd.Dir = cwd + "/smart-contracts/"
	_, err = npmInstWeb3Cmd.Output()
	if err != nil {
		log.Println(err)
	}
	// fmt.Printf("%s", output)

	npmInstSolcCmd := exec.Command("npm", "install", "solc@0.4.25")
	npmInstSolcCmd.Dir = cwd + "/smart-contracts/"
	_, err = npmInstSolcCmd.Output()
	if err != nil {
		log.Println(err)
	}
	// fmt.Printf("%s", output)

	fmt.Println("\rDependencies has been successfully generated.")
}

func deployContract(fileName, IP string) string {
	fmt.Println("Deploying Smart Contract: " + fileName)
	cwd := os.Getenv("HOME")
	deployCmd := exec.Command("node", "deploy.js", fileName, IP)
	deployCmd.Dir = cwd + "/smart-contracts/"
	output, err := deployCmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", output)
	return fmt.Sprintf("%s", output)
}

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

var gethSolcCmd = &cobra.Command{
	Use:   "solc",
	Short: "Smart contract deployment tool",
	Long: `
Solc will allow the user to reploy smart contracts to the ethereum blockchain.
	`,
	Run: util.PartialCommand,
}

var gethSolcInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the smart-contracts directory",
	Long: `
Init initialize the smart-contracts directory and will download all the necessary dependencies. This may take some time as the files are being pulled. All smart contracts should be put into the 'smart-contracts' directory.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking Directory")
		checkContractDir()
		installNpmDeps()
		fmt.Println("'smart-contracts' directory is initializd and smart contract deployment is now available.")
	},
}

var gethSolcDeployCmd = &cobra.Command{
	Use:   "deploy <node> <file name>",
	Short: "deploy",
	Long: `
Deploy will compile the smart contract and deploy it to the ethereum blockchain. For the smart contract to be successfully deployed, mining needs to be started. This can be done by using the 'miner start' command. 

Output: Deployed contract address
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 2, 2)
		if checkContractFiles(args[1]) {
			nodes, err := GetNodes()
			if err != nil {
				util.PrintErrorFatal(err)
			}

			nodeNumber, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}

			nodeIP := nodes[nodeNumber].IP
			deployContractOut := deployContract(args[1], nodeIP)
			re := regexp.MustCompile(`(?m)0x[0-9a-fA-F]{40}`)
			addrList := re.FindAllString(deployContractOut, -1)
			// fmt.Println(addrList[:])

			ContractList := make([]interface{}, 0)
			ContractList = append(ContractList, Contracts{
				DeployedNodeAddress: addrList[0],
				ContractName:        args[1],
				ContractAddress:     addrList[1],
			})
			contracts, _ := json.Marshal(ContractList)
			writeContractListFile(fmt.Sprintf("%s", contracts))
		}
	},
}

var gethConsole = &cobra.Command{
	Use:   "console <node>",
	Short: "Logs into the geth console",
	Long: `
Console will log into the geth console.

Response: stdout of geth console`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		nodes, err := GetNodes()
		if err != nil {
			util.PrintErrorFatal(err)
		}

		nodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}

		log.Fatal(unix.Exec("/usr/bin/ssh", []string{"ssh", "-i", "/home/master-secrets/id.master", "-o", "StrictHostKeyChecking no",
			"-o", "UserKnownHostsFile=/dev/null", "-o", "PasswordAuthentication no", "-y",
			"root@" + fmt.Sprintf(nodes[nodeNumber].IP), "-t", "tmux", "attach", "-t", "whiteblock"}, os.Environ()))
	},
}

var gethGetTxReceiptCmd = &cobra.Command{
	Use:   "get_transaction_receipt <hash>",
	Short: "Get transaction receipt",
	Long: `
Get the transaction receipt by the tx hash

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction receipt.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("eth::get_transaction_receipt", args)
	},
}

var gethGetHashRateCmd = &cobra.Command{
	Use:   "get_hash_rate",
	Short: "Get hash rate",
	Long: `
Get the current hash rate per node

Response: The hash rate of a single node in the network`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 0, 1)
		jsonRpcCallAndPrint("eth::get_hash_rate", []string{})
	},
}

var gethGetRecentSentTxCmd = &cobra.Command{
	Use:   "get_recent_sent_tx <number>",
	Short: "Get recently sent transaction",
	Long: `
Get a number of the most recent transactions sent

Format: <number>
Params: The number of transactions to retrieve

Response: JSON object of transaction data`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("eth::get_recent_sent_tx", args)
	},
}

/*
var gethGetBlockNumberCmd = &cobra.Command{
	Use:   "get_block_number",
	Short: "Get block number",
	Long: `
Get the current highest block number of the chain

Response: The block number`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 0, 1)
		jsonRpcCallAndPrint("eth::get_block_number", []string{})
	},
}

var gethGetBlockCmd = &cobra.Command{
	Use:   "get_block <block number>",
	Short: "Get block information",
	Long: `
Get the data of a block

Format: <Block Number>
Params: Block number

Response: JSON Representation of the block.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("eth::get_block", args)
	},
}

var gethGetAccountCmd = &cobra.Command{
	Use:   "get_accounts",
	Short: "Get account information",
	Long: `
Get a list of all unlocked accounts, current balance of accounts, tx counts, and other relevant information.

Response: A JSON array of the accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("eth::accounts_status", []string{})
	},
}

var gethSendTxCmd = &cobra.Command{
	Use:   "send_transaction <from address> <to address> <gas> <gas price> <value to send>",
	Short: "Sends a transaction",
	Long: `
Send a transaction between two accounts

Format: <from> <to> <gas> <gas price> <value>
Params: Sending account, receiving account, gas, gas price, amount to send in ETH

Response: The transaction hash`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 5, 5)
		args[4] = args[4] + "000000000000000000"
		jsonRpcCallAndPrint("eth::send_transaction", args)
	},
}

var gethGetTxCountCmd = &cobra.Command{
	Use:   "get_transaction_count <address> [block number]",
	Short: "Get transaction count",
	Long: `
Get the transaction count sent from an address, optionally by block

Format: <address> [block number]
Params: The sender account, a block number

Response: The transaction count`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 2)
		jsonRpcCallAndPrint("eth::get_transaction_count", args)
	},
}

var gethGetTxCmd = &cobra.Command{
	Use:   "get_transaction <hash>",
	Short: "Get transaction information",
	Long: `
Get a transaction by its hash

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		jsonRpcCallAndPrint("eth::get_transaction", args)
	},
}

var gethStartTxCmd = &cobra.Command{
	Use:   "start_transactions <tx/s> <value> [destination]",
	Short: "Start transactions",
	Long: `
Start sending transactions according to the given parameters.

Format: <tx/s> <value> [destination]
Params: The amount of transactions to send in a second, the value of each transaction in wei, the destination for the transaction
`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(command)
		util.CheckArguments(args, 2, 3)
		args[1] = args[1] + "000000000000000000"
		jsonRpcCallAndPrint("eth::start_transactions", args)
	},
}

var gethStopTxCmd = &cobra.Command{
	Use:   "stop_transactions",
	Short: "Stop transactions",
	Long: `
Stops the sending of transactions if transactions are currently being sent`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 0, 1)
		jsonRpcCallAndPrint("eth::stop_transactions", []string{})
	},
}

var gethStartMiningCmd = &cobra.Command{
	Use:   "start_mining [node 1 number] [node 2 number]...",
	Short: "Start Mining",
	Long: `
Send the start mining signal to nodes, may take a while to take effect due to DAG generation

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to start mining or None for all nodes

Response: The number of nodes which successfully received the signal to start mining`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("eth::start_mining", args)
		DagReady := false
		for !DagReady {
			fmt.Printf("\rDAG is being generated...")
			res, err := jsonRpcCall("eth::get_block_number", []string{})
			if err != nil {
				util.PrintErrorFatal(err)
			}
			blocknum := res.(int)
			if blocknum > 4 {
				DagReady = true
			}
		}
		fmt.Println("\rDAG has been successfully generated.")
	},
}

var gethStopMiningCmd = &cobra.Command{
	Use:   "stop_mining [node 1 number] [node 2 number]...",
	Short: "Stop mining",
	Long: `
Send the stop mining signal to nodes

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to stop mining or None for all nodes


Response: The number of nodes which successfully received the signal to stop mining`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("eth::stop_mining", args)
	},
}

*/

func init() {
	// gethCmd.Flags().StringVarP(&gethcommand, "command", "c", "", "Geth command")
	gethCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	//geth subcommands
	gethCmd.AddCommand(gethGetTxReceiptCmd, gethGetHashRateCmd,
		gethGetRecentSentTxCmd, gethConsole, gethSolcCmd)

	gethSolcCmd.AddCommand(gethSolcInitCmd, gethSolcDeployCmd)
	RootCmd.AddCommand(gethCmd)
}
