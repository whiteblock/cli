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
	Long: `
Build will create and deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own container and will interact individually as a participant of the specified network.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		command := "build"
		param := "{\"servers\":" + fmt.Sprintf("%s", server) + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + fmt.Sprintf("%d", nodes) + ",\"image\":\"" + image + "\"}"
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		wsEmitListen(serverAddr, command, param)

		for _, serv := range server {
			switch serv {
			case "1":
				println("GUI to view stats and network information found here: 172.16.1.5:3000")
			case "2":
				println("GUI to view stats and network information found here: 172.16.2.5:3000")
			case "3":
				println("GUI to view stats and network information found here: 172.16.3.5:3000")
			case "4":
				println("GUI to view stats and network information found here: 172.16.4.5:3000")
			case "5":
				println("GUI to view stats and network information found here: 172.16.5.5:3000")
			}
		}
	},
}

func init() {
	buildCmd.Flags().StringVarP(&blockchain, "blockc", "b", "ethereum", "blockchain")
	buildCmd.Flags().StringVarP(&image, "image", "i", "ethereum:latest", "image")
	buildCmd.Flags().IntVarP(&nodes, "nodes", "n", 5, "number of nodes")
	buildCmd.Flags().StringArrayVarP(&server, "server", "s", []string{}, "servers to build on")
	buildCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(buildCmd)
}
