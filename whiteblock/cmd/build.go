package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	blockchain string
	image      string
	nodes      int
	server     []string
	serverAddr string
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: `Build will create and deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own containers and will interact individually as a participant of the specified network.
	
	image 				Specifies the docker image of the network to deploy, default 'Geth' image will be used;
	nodes 				Number of nodes to create, 10 will be used as default;
	server 				Number of servers to deploy network, e.g. "alpha", "bravo", "charlie", etc.;
	`,

	Run: func(cmd *cobra.Command, args []string) {
		// curlPOST(fmt.Sprint(serverAddr)+"/testnets/", "-d '{\"Servers\":\""+fmt.Sprint("%d", server)+"\",\"Blockchain\":\"ethereum\",\"Nodes\":"+fmt.Sprintf("%d", nodes)+",\"Image\":\""+fmt.Sprint(image)+"\"}'")
		msg := "build,{\"Servers\":" + fmt.Sprintf("%s", server) + ",\"Blockchain\":\"" + blockchain + "\",\"Nodes\":" + fmt.Sprintf("%d", nodes) + ",\"Image\":\"" + image + "\"}"
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		wsBuild(serverAddr, msg)

	},
}

func init() {
	buildCmd.Flags().StringVarP(&blockchain, "blockc", "b", "ethereum", "blockchain")
	buildCmd.Flags().StringVarP(&image, "image", "i", "ethereum:latest", "image")
	buildCmd.Flags().IntVarP(&nodes, "nodes", "n", 10, "number of nodes")
	buildCmd.Flags().StringArrayVarP(&server, "server", "s", []string{}, "number of servers")
	buildCmd.Flags().StringVarP(&serverAddr, "serverAddr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(buildCmd)
}
