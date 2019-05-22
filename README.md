# Whiteblock CLI
[![Maintainability](https://api.codeclimate.com/v1/badges/19632596f75488519c67/maintainability)](https://codeclimate.com/github/Whiteblock/cli/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/whiteblock/cli)](https://goreportcard.com/report/github.com/whiteblock/cli)
* Latest Stable Build (Linux AMD64)

  https://storage.cloud.google.com/genesis-public/cli/master/bin/linux/amd64/whiteblock

* Latest Dev Build (Linux AMD64)

  https://storage.cloud.google.com/genesis-public/cli/dev/bin/linux/amd64/whiteblock

## ./whiteblock <COMMAND> [FLAGS]
This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation, usages, and exmaples can be found in our [documentation](www.whiteblock.io/docs/cli).

* Available Commands:
    * build       Build a blockchain using image and deploy nodes
    * get         Get server and network information.
    * geth        Run geth commands
    * help        Help about any command
    * netconfig   Network conditions
    * rpc         Rpc interacts with the blockchain
    * ssh         SSH into an existing container.
    * version     Get whiteblock CLI client version

* Flags:
  *  -h, --help : help for whiteblock

### build
./whiteblock build [FLAGS]

Aliases: build, create, init

Build will create and deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own container and will interact individually as a participant of the specified network.

* Flags:
  *  -h, --help:                 help for build
  *  -a, --server-addr string:   server address with port 5000 (default "localhost:5000")

### get
./whiteblock get <command> [FLAGS]

Get will ouput server and network information and statstics.

* Available Commands:
    * data        Data will pull data from the network and output into a file.
    * nodes       Nodes will show all nodes in the network.
    * server      Get server information.
    * stats       Get stastics of a blockchain

* Flags:
  *  -h, --help : help for get
  *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

### get data
./whiteblock get data <command> [FLAGS]

Data will pull specific or all block data from the network and output into a file. You will specify the directory where the file will be downloaded.

* Available Commands:
    * all         All will pull data from the network and output into a file.
    * block       Data block will pull data from the network and output into a file.
    * time        Data time will pull data from the network and output into a file.

* Flags:
    *  -h, --help : help for data
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get data all
./whiteblock get data all [path] [FLAGS]

Data all will pull all data from the network and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

Response: JSON representation of network statistics

* Flags:
  *  -h, --help : help for all
  *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get data block
./whiteblock get data block <start block> <end block> [path] [FLAGS]

Data block will pull block data from the network from a given start and end block and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

Params: Block numbers
Format: <start block number> <end block number>

Response: JSON representation of network statistics

* Flags:
  *  -h, --help : help for block
  *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get data time
./whiteblock get data time <start time> <end time> [path] [FLAGS]

Data time will pull block data from the network from a given start and end time and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

Response: JSON representation of network statistics

* Flags:
  *  -h, --help : help for time
  *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get default
./whiteblock get default <blockchain> [FLAGS]

Get the blockchain specific parameters for a deployed blockchain.

Params: The blockchain to get the build params of
Format: <blockchain>

Response: The params as a list of key value params, of name and type respectively

* Flags:
  *  -h, --help : help for default

#### get nodes
./whiteblock get nodes [FLAGS]

Aliases: nodes, node

Nodes will output all of the nodes in the current network.

* Flags:
    *  -h, --help : help for server
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get server
./whiteblock get server [FLAGS]

Aliases: server, servers

Server will allow the user to get server information.

* Flags:
    *  -h, --help : help for server
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

### get stats
./whiteblock get stats <command> [FLAGS]

Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics

* Available Commands:
    * all
    * block
    * time

* Flags:
    *  -h, --help : help for stats
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get stats all
./whiteblock get stats all [FLAGS]

Stats all will allow the user to get all the statistics regarding the network.

Response: JSON representation of network statistics

* Flags:
    *  -h, --help : help for all
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get stats block
./whiteblock get stats block <start block> <end block> [FLAGS]

Stats block will allow the user to get statistics regarding the network.

Params: Block numbers
Format: <start block number> <end block number>

Response: JSON representation of network statistics

* Flags:
    *  -h, --help : help for block
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get stats time
./whiteblock get stats time <start time> <end time> [FLAGS]
Stats time will allow the user to get statistics by specifying a start time and stop time (unix time stamp).

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

Response: JSON representation of network statistics

* Flags:
    *  -h, --help : help for time
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

### geth
./whiteblock geth <command> [flags]

Geth will allow the user to get infromation and run geth commands.

* Available SubCommands:
    * block_listener          Get block listener
    * get_accounts            Get account information
    * get_balance             Get account balance information
    * get_block               Get block information
    * get_block_number        Get block number
    * get_hash_rate           Get hasg rate
    * get_recent_sent_tx      Get recently sent transaction
    * get_transaction         Get transaction information
    * get_transaction_count   Get transaction count
    * get_transaction_receipt Get transaction receipt
    * send_transaction        Sends a transaction
    * start_mining            Start Mining
    * start_transactions      Start transactions
    * stop_mining             Stop mining
    * stop_transactions       Stop transactions

* Flags:
    *  -h, --help:               help for geth
    *  -a, --server-addr `string`:   server address with port 5000 (default "localhost:5000")

#### geth block_listener
./whiteblock geth block_listener [block number] [flags]

Get all blocks and continue to subscribe to new blocks

Format: [block number]
Params: The block number to start at or None for all blocks
Response: Will emit on eth::block_listener for every block after the given block or 0 that exists/has been created

* Flags:
    *  -h, --help:   help for block_listener

#### geth get_accounts
./whiteblock geth get_accounts [flags]

Get a list of all unlocked accounts

Response: A JSON array of the accounts

* Flags:
    * -h, --help:   help for get_accounts

#### geth get_balance
./whiteblock geth get_balance <address> [flags]

Get the current balance of an account

Format: <address>
Params: Account address
Response: The integer balance of the account in wei

* Flags:
    *  -h, --help:   help for get_balance

#### geth get_block
./whiteblock geth get_block <block number> [flags]

Get the data of a block

Format: <Block Number>
Params: Block number

* Flags:
    * -h, --help:   help for get_block

#### geth get_block_number
./whiteblock geth get_block_number [flags]

Get the current highest block number of the chain

Response: The block number

* Flags:
    *  -h, --help:   help for get_block_number

#### geth get_hash_rate
./whiteblock geth get_hash_rate [flags]

Get the current hash rate per node

Response: The hash rate of a single node in the network

* Flags:
    *  -h, --help:   help for get_hash_rate
  
#### geth get_recent_sent_tx
./whiteblock geth get_recent_sent_tx [number] [flags]

Get a number of the most recent transactions sent

Format: [number]
Params: The number of transactions to retrieve
Response: JSON object of transaction data

* Flags:
    *  -h, --help:   help for get_recent_sent_tx

#### geth get_transaction
./whiteblock geth get_transaction <hash> [flags]

Get a transaction by its hash

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction.

* Flags:
    *  -h, --help:   help for get_transaction

#### geth get_transaction_count
./whiteblock geth get_transaction_count <address> [block number] [flags]

Get the transaction count sent from an address, optionally by block

Format: <address> [block number]
Params: The sender account, a block number
Response: The transaction count

* Flags:
    *  -h, --help:   help for get_transaction_count

#### geth get_transaction_receipt
./whiteblock geth get_transaction_receipt <hash> [flags]

Get the transaction receipt by the tx hash

Format: <hash>
Params: The transaction hash
Response: JSON representation of the transaction receipt.

* Flags:
    *  -h, --help:   help for get_transaction_receipt

#### geth send_transaction
./whiteblock geth send_transaction <from address> <to address> <gas> <gas price> <value to send> [flags]

Send a transaction between two accounts

Format: <from> <to> <gas> <gas price> <value>
Params: Sending account, receiving account, gas, gas price, amount to send, transaction data, nonce
Response: The transaction hash

* Flags:
    *  -h, --help:   help for send_transaction

#### geth start_mining
./whiteblock geth start_mining [node 1 number] [node 2 number]... [flags]

Send the start mining signal to nodes, may take a while to take effect due to DAG generation

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to start mining or None for all nodes
Response: The number of nodes which successfully received the signal to start mining

* Flags:
    *  -h, --help:   help for start_mining

#### geth start_transactions
./whiteblock geth start_transactions <tx/s> <value> [destination] [flags]

Start sending transactions according to the given parameters, value = -1 means randomize value.

Format: <tx/s> <value> [destination]
Params: The amount of transactions to send in a second, the value of each transaction in wei, the destination for the transaction

* Flags:
    *  -h, --help:   help for start_transactions

#### geth stop_mining
./whiteblock geth stop_mining [node 1 number] [node 2 number]... [flags]

Send the stop mining signal to nodes

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to stop mining or None for all nodes
Response: The number of nodes which successfully received the signal to stop mining

* Flags:
    *  -h, --help:   help for stop_mining

#### geth stop_transactions
./whiteblock geth stop_transactions [flags]

Stops the sending of transactions if transactions are currently being sent

* Flags:
    *  -h, --help:   help for stop_transactions

### netconfig
./whiteblock netconfig <command> [FLAGS]

Aliases: emulate

Netconfig will introduce persisting network conditions for testing. Use '?' at any time for more help on configuring the network.

Custom Command:
netconfig <engine number> <path number> <command>

set delay <amount> 			Specifies the latency to add [ms];
set loss loss <amount>			Specifies the amount of packet loss to add [%];
set bw <amount> <type>			Specifies the bandwidth of the network [bps|Kbps|Mbps|Gbps];

* Available Commands:
    * bandwidth   Set bandwidth
    * delay       Set latency
    * loss        Set packetloss
    * off         Turn off emulation
    * on          Turn on emulation

* Flags:
    *  -h, --help:   help for netconfig

#### netconfig bandwidth
./whiteblock netconfig bandwidth <engine number> <path number> <amount> <bandwidth type> [FLAGS]

Aliases: bw

Bandwidth will constrict the network to the specified bandwidth. You will specify the amount of bandwdth and the type.

Fomat: 
	bandwidth type: bps, Kbps, Mbps, Gbps

* Flags:
    *  -h, --help:   help for bandwidth

#### netconfig delay
./whiteblock netconfig delay <engine number> <path number> <amount> [FLAGS]

Aliases: delay, latancy, lat

Latency will introduce delay to the network. You will specify the amount of latency in ms.

* Flags:
    *  -h, --help:   help for latency

#### netconfig loss
./whiteblock netconfig loss <engine number> <path number> <percent> [FLAGS]

Aliases: packetloss

Packetloss will drop packets in the network. You will specify the amount of packet loss in %.

* Flags:
    *  -h, --help:   help for loss

#### netconfig off
./whiteblock netconfig off <engine number> [FLAGS]

Turn off emulation.

* Flags:
  *  -h, --help:   help for off

#### netconfig on
./whiteblock netconfig on <engine number> [FLAGS]

Turn on emulation.

* Flags:
    *  -h, --help:   help for on

### ssh
./whiteblock ssh <server> <node> [FLAGS]

SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command

* Flags:
  *  -h, --help :          help for ssh
  *  -a, --server-addr :   server address with port 5000 (default "localhost:5000")

### sys
./whiteblock sys <command> [FLAGS]

Alias: SYS, syscoin

Sys will allow the user to get infromation and run SYS commands.

* Available Commands:
  * test        SYS test commands.

* Flags:
  *  -h, --help :          help for sys

#### sys test
./whiteblock sys test <command> [FLAGS]

Available Commands:
  results     Get results from a previous test.
  start       Starts propagation test.

* Flags:
  *  -h, --help :          help for test

#### sys test start
./whiteblock sys test start <wait time> <min complete percent> <number of tx> [FLAGS]

Sys test start will start the propagation test. It will wait for the signal start time, have nodes send messages at the same time, and require to wait a minimum amount of time then check receivers with a completion rate of minimum completion percentage.

Format: <wait time> <min complete percent> <number of tx>
Params: Time in seconds, percentage, number of transactions

* Flags:
    *  -h, --help : help for start
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### sys test results
./whiteblock sys test results <test number> [FLAGS]

Sys test results pulls data from a previous test or tests and outputs as csv.

Format: <test number>
Params: Test number

* Flags:
    *  -h, --help : help for results
    *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")


### version [FLAGS]
./whiteblock version

Get whiteblock CLI client version

* Flags:
  *  -h, --help : help for version


**** TO CONFIGURE: ****

### contractadd
./whiteblock contractadd <filename> [flags]

Adds the specified smart contract into the /Downloads folder.

* Flags:
  *  -h, --help:              help for contractadd
  *  -p, --path `string` :      File path where the smart contract is located

### contractcompile
./whiteblock contractcompile <filename> [flags]

Compiles the specified smart contract.

* Flags:
  * -h, --help:              help for contractcompile
  * -p, --path `string`:       File path where the smart contract is located


