package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	image      string
	nodes      int
	server     int
	serverAddr string
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: `Build will create and deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own containers and will interact individually as a participant of the specified network.
	
	image 				Specifies the docker image of the network to deploy, default 'Geth' image will be used;
	nodes 				Number of nodes to create, 10 will be used as default;
	server 				Number of servers to deploy network, 1 server will be used;
	`,

	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("http://localhost:8000/testnets/", "-d '{\"Servers\":\""+fmt.Sprint(server)+"\",\"Blockchain\":\"ethereum\",\"Nodes\":"+fmt.Sprintf("%d", nodes)+",\"Image\":\"ethereum:latest\"}'")
		// fmt.Println(image, nodes, server)

		curlPOST(fmt.Sprint(serverAddr)+"/testnets/", "-d '{\"Servers\":\""+fmt.Sprint("%d", server)+"\",\"Blockchain\":\"ethereum\",\"Nodes\":"+fmt.Sprintf("%d", nodes)+",\"Image\":\""+fmt.Sprint(image)+"\"}'")
	},
}

func init() {
	buildCmd.Flags().StringVarP(&image, "image", "i", "ethereum:latest", "image")
	buildCmd.Flags().IntVarP(&nodes, "nodes", "n", 10, "number of nodes")
	buildCmd.Flags().IntVarP(&server, "server", "s", 1, "number of servers")
	buildCmd.Flags().StringVarP(&serverAddr, "serverAddr", "a", "http://localhost:8000", "server address with port 8000")

	RootCmd.AddCommand(buildCmd)
}
