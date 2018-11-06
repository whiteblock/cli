# Whiteblock CLI

./whiteblock [COMMAND] [FLAGS]

This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation and exmaples can be found at www.whiteblock.io/docs/cli.
* flags:
    * -h, --help : help for whiteblock

## Commands

### init
./whiteblock init [FLAGS]

* aliases: build, create, init
Build will deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own containers and will interact individually as a p
articipant of the specified blockchain.

    Image                 Specifies the docker image of the network to deploy, default 'Geth' image will be used;
    Nodes                Number of nodes to create, 10 will be used as default;
    Server                Number of servers to deploy network, 1 server will be used;

* flags:
    * -h, --help : help for init
    * -i, --image string : image (default "ethereum:latest")
    * -n, --nodes int : number of nodes (default 10)
    * -s, —server int : number of servers (default 1)

### send
./whiteblock send [FLAGS]

Send will have nodes send a specified number of transactions from every node that had been deployed.

        transactions              Sends specified number of transactions;
        senders                       Number of nodes sending transactions;

* flags:
    * -t, --transactions int : number of transactions to send (default 100)
    * -s, --senders int : Number of Nodes Sending (default 10)

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