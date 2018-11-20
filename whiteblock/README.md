# Whiteblock CLI

./whiteblock <COMMAND> [FLAGS]

This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation, usages, and exmaples can be found at [www.whiteblock.io/docs/cli].

* Available Commands:
    * build           Build a blockchain using image and deploy nodes
    * contractadd     Add a smart contract.
    * contractcompile Smart contract compiler.
    * get             Get server and network information.
    * geth            Run geth commands
    * help            Help about any command
    * netconfig       Network conditions
    * send            Send transactions from all nodes
    * ssh             SSH into an existing container.
    * version         Get whiteblock CLI client version

* flags:
    * -h, --help : help for whiteblock

## Commands

### build / init / create
./whiteblock build [FLAGS]

* aliases: build, create, init
Build will deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own containers and will interact individually as a p
articipant of the specified blockchain.

* flags:
    * -b, --blockc string        blockchain (default "ethereum")
    * -h, --help                 help for build
    * -i, --image string         image (default "ethereum:latest")
    * -n, --nodes int            number of nodes (default 10)
    * -s, --server stringArray   number of servers
    * -a, --server-addr string   server address with port 5000 (default "localhost:5000")

### get
./whiteblock get <COMMAND> [FLAGS]

Get will allow the user to get server and network information.

* flags:
  -h, --help                 help for get
  -a, --server-addr string   server address with port 5000 (default "localhost:5000")

### server
./whiteblock server [FLAGS]

Server will allow the user to get server information.

        serverInfo                                       Get the information from all currently registered servers;
        serverInfo --id [Server ID]             Get server information by id;

* flags:
    * -i, --ID string : server ID
    * -h, --help : help for server

### testnet
./whiteblock testnet [FLAGS]

Testnet will allow the user to get infromation regarding the test network.

        testnetInfo                                       Get all testnets which are currently running
        testnetInfo --id [Testnet ID]           Get data on a single testnet
        addTestnet                                      Add and deploy a new testnet

* flags:
    * -i, --ID string : testnet ID
    * -h, --help : help for testnet


### version 
./whiteblock version

Get whiteblock CLI client version

* flags:
    * -h, --help : help for version

### netropy
./whiteblock netrop [FLAGS]

Netropy will introduce persisting network conditions for testing.

        latency                         Specifies the latency to add [ms];
        packetloss                      Specifies the amount of packet loss to add [%];

* flags:
    * -l, --latency int : latency (default 10)
    *   -p, --packetloss float : packetloss (default 0.001)