# Whiteblock CLI

## ./whiteblock <COMMAND> [FLAGS]
This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation, usages, and exmaples can be found at [www.whiteblock.io/docs/cli].

* Available Commands:
    * build           
    * contractadd     
    * contractcompile 
    * get             
    * geth
    * help 
    * netconfig
    * send    
    * ssh 
    * version 

* Flags:
  *  -h, --help : help for whiteblock

### build
./whiteblock build [FLAGS]

Aliases: build, create, init

Build will deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own containers and will interact individually as a participant of the specified blockchain.

* Flags:
  *  -b, --blockc string:        blockchain (default "ethereum")
  *  -h, --help:                 help for build
  *  -i, --image string:         image (default "ethereum:latest")
  *  -n, --nodes int:            number of nodes (default 10)
  *  -s, --server stringArray:   number of servers
  *  -a, --server-addr string:   server address with port 5000 (default "localhost:5000")

### get <subcommand>
./whiteblock get <SUBCOMMAND> [FLAGS]

Get will allow the user to get server and network information.

* Available Commands:
    * nodes
    * server
    * testnet

* Flags:
  *  -h, --help : help for get
  *  -a, --server-addr `string`:  server address with port 5000 (default "localhost:5000")

#### get nodes [FLAGS]
./whiteblock get nodes [FLAGS]

Aliases: nodes, node

Nodes will output all of the nodes in the current network.

* Flags:
  *  -h, --help : help for server

#### get server [FLAGS]
./whiteblock get server [FLAGS]

Aliases: server, servers

Server will allow the user to get server information.

* Flags:
  *  -h, --help : help for server

#### get testnet [FLAGS]
./whiteblock get testnet [FLAGS]

Testnet will allow the user to get infromation regarding the test network.

* Flags:
  *  -h, --help : help for testnet

### geth <SUBCOMMAND> [FLAGS]
./whiteblock geth <subcommand> [flags]

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

#### geth block_listener [block number] [FLAGS]
./whiteblock geth block_listener [block number] [flags]

Get all blocks and continue to subscribe to new blocks

Format: [block number]
Params: The block number to start at or None for all blocks
Response: Will emit on eth::block_listener for every block after the given block or 0 that exists/has been created

* Flags:
  *  -h, --help:   help for block_listener

#### geth get_accounts [FLAGS]
./whiteblock geth get_accounts [flags]

Get a list of all unlocked accounts

Response: A JSON array of the accounts

* Flags:
  * -h, --help:   help for get_accounts

#### geth get_balance <ADDRESS> [FLAGS]
./whiteblock geth get_balance <address> [flags]

Get the current balance of an account

Format: <address>
Params: Account address
Response: The integer balance of the account in wei

* Flags:
  *  -h, --help:   help for get_balance

#### geth get_block <BLOCK NUMBER> [FLAGS]
./whiteblock geth get_block <block number> [flags]

Get the data of a block

Format: <Block Number>
Params: Block number

* Flags:
  * -h, --help:   help for get_block

#### geth get_block_number [FLAGS]
./whiteblock geth get_block_number [flags]

Get the current highest block number of the chain

Response: The block number

* Flags:
  *  -h, --help:   help for get_block_number

#### geth get_hash_rate [FLAGS]
./whiteblock geth get_hash_rate [flags]

Get the current hash rate per node

Response: The hash rate of a single node in the network

* Flags:
  *  -h, --help:   help for get_hash_rate
  
#### geth get_recent_sent_tx [NUMBER] [FLAGS]
./whiteblock geth get_recent_sent_tx [number] [flags]

Get a number of the most recent transactions sent

Format: [number]
Params: The number of transactions to retrieve
Response: JSON object of transaction data

* Flags:
  *  -h, --help:   help for get_recent_sent_tx

#### geth get_transaction <HASH> [FLAGS]
./whiteblock geth get_transaction <hash> [flags]

Get a transaction by its hash

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction.

* Flags:
  *  -h, --help:   help for get_transaction

#### geth get_transaction_count <ADDRESS> [BLOCK NUMBER [FLAGS]
./whiteblock geth get_transaction_count <address> [block number] [flags]

Get the transaction count sent from an address, optionally by block

Format: <address> [block number]
Params: The sender account, a block number
Response: The transaction count

* Flags:
  *  -h, --help:   help for get_transaction_count

#### geth get_transaction_receipt <HASH> [FLAGS]
./whiteblock geth get_transaction_receipt <hash> [flags]

Get the transaction receipt by the tx hash

Format: <hash>
Params: The transaction hash
Response: JSON representation of the transaction receipt.

* Flags:
  *  -h, --help:   help for get_transaction_receipt

#### geth send_transaction <FROM> <TO> <GAS> <GAS PRICE> <VALUE> [FLAGS]
./whiteblock geth send_transaction <from address> <to address> <gas> <gas price> <value to send> [flags]

Send a transaction between two accounts

Format: <from> <to> <gas> <gas price> <value>
Params: Sending account, receiving account, gas, gas price, amount to send, transaction data, nonce
Response: The transaction hash

* Flags:
  *  -h, --help:   help for send_transaction

#### geth start_mining [NODE 1 NUMBER] [NODE 2 NUMBER] ... [FLAGS]
./whiteblock geth start_mining [node 1 number] [node 2 number]... [flags]

Send the start mining signal to nodes, may take a while to take effect due to DAG generation

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to start mining or None for all nodes
Response: The number of nodes which successfully received the signal to start mining

* Flags:
  *  -h, --help:   help for start_mining

#### geth start_transactions <TX/S> <VALUE> [DESTINATION] [FLAGS]
./whiteblock geth start_transactions <tx/s> <value> [destination] [flags]

Start sending transactions according to the given parameters, value = -1 means randomize value.

Format: <tx/s> <value> [destination]
Params: The amount of transactions to send in a second, the value of each transaction in wei, the destination for the transaction

* Flags:
  *  -h, --help:   help for start_transactions

#### geth stop_mining [NODE 1 NUMBER] [NODE 2 NUMBER] ... [FLAGS]
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

### ssh [FLAGS]
./whiteblock ssh [flags]

SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command

* Flags:
  *  -c, --cmd `string`:           Which shell to run in container (default "ls -l")
  *  -h, --help:                 help for ssh
  *  -n, --node `int`:             Node number to SSH into
  *  -s, --server `string`:        Which server to run in (default "1")
  *  -a, --server-addr `string`:   server address with port 5000 (default "localhost:5000")

### version [FLAGS]
./whiteblock version

Get whiteblock CLI client version

* Flags:
  *  -h, --help : help for version


**** TO CONFIGURE: ****


### netconfig <SUBCOMMAND> [FLAGS]
./whiteblock netconfig [COMMAND] [ENGINE NUMBER] [ENGINE NUMBER] [FLAGS]

Netropy will introduce persisting network conditions for testing.

* Available Commands:
    * latency           
    * packetloss     
    * bandwidth 

#### netropy latency <MILISEC> [FLAGS]
./whiteblock netconfig latency <milisec> [flags]

Latency will set persisting latency to the network.

#### netropy packetloss <PERCENT> [FLAGS]
./whiteblock netconfig packetloss <percent> [flags]

Packetloss will drop random packets in the network.

#### netropy bandwidth <MB> [FLAGS]
./whiteblock netconfig bandwidth <mb> [flags]

Bandwidth will delegate a specific amount of bandwidth to the network.


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


