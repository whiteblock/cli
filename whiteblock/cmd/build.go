package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	serverAddr string

	// blockchain string
	// image      string
	// nodes      int
	// server     []string
	// cpu        string
	// memory     string
	// params     string
)

func checkServ(server string) string {
	servList := make([]string, 0)
	servList = append(servList, "GUI to view stats and network information found here:")
	if strings.Contains(server, "1") {
		servList = append(servList, " 172.16.1.5:3000")
	}
	if strings.Contains(server, "2") {
		servList = append(servList, " 172.16.2.5:3000")
	}
	if strings.Contains(server, "3") {
		servList = append(servList, " 172.16.3.5:3000")
	}
	if strings.Contains(server, "4") {
		servList = append(servList, " 172.16.4.5:3000")
	}
	if strings.Contains(server, "5") {
		servList = append(servList, " 172.16.5.5:3000")
	}
	return strings.Join(servList, " ")
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: `
Build will create and deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own container and will interact individually as a participant of the specified network.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		buildArr := make([]string, 0)
		paramArr := make([]string, 0)
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "build"

		buildOpt := [6]string{"servers (default set to: [])", "blockchain (default set to: ethereum)", "nodes (default set to: 10)", "image (default set to: ethereum:latest)", "cpus (default set to: \"\")", "memory (default set to: \"\")"}
		defOpt := [6]string{"[]", "ethereum", "10", "ethereum:latest", "", ""}

		scanner := bufio.NewScanner(os.Stdin)
		for i := 0; i < len(buildOpt); i++ {
			fmt.Print(buildOpt[i] + ": ")
			scanner.Scan()

			text := scanner.Text()
			if len(text) != 0 {
				buildArr = append(buildArr, text)
			} else {
				buildArr = append(buildArr, defOpt[i])
			}
		}

		server := "[" + buildArr[0] + "]"
		blockchain := buildArr[1]
		nodes := buildArr[2]
		image := buildArr[3]
		cpu := buildArr[4]
		memory := buildArr[5]

		// param := "{\"servers\":" + fmt.Sprintf("%s", server) + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"}}"

		if blockchain == "ethereum" {
			q := [9]string{"chainId", "networkId", "difficulty", "initBalance", "maxPeers", "gasLimit", "homesteadBlock", "eip155Block", "eip158Block"}

			scanner := bufio.NewScanner(os.Stdin)
			for i := 0; i < len(q); i++ {
				fmt.Print(q[i] + ": ")
				scanner.Scan()

				text := scanner.Text()
				if len(text) != 0 {
					paramArr = append(paramArr, text)
				} else {
					paramArr = append(paramArr, "\"\"")
				}
			}

			bcParameters := "{\"chainId\":" + paramArr[0] + ",\"networkId\":" + paramArr[1] + ",\"difficulty\":" + paramArr[2] + ",\"initBalance\":" + paramArr[3] + ",\"maxPeers\":" + paramArr[4] + ",\"gasLimit\":" + paramArr[5] + ",\"homesteadBlock\":" + paramArr[6] + ",\"eip155Block\":" + paramArr[7] + ",\"eip158Block\":" + paramArr[8] + "}"

			param := "{\"servers\":" + fmt.Sprintf("%s", server) + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"},\"params\":" + bcParameters + "}"
			wsEmitListen(serverAddr, command, param)

		} else if blockchain == "syscoin" {
			q := [4]string{"rpcUser", "rpcPass", "difficulty", "extras"}
			o := [9]string{"server", "regtest", "listen", "rest", "debug", "unittest", "addressindex", "assetallocationindex", "tpstest"}

			scanner := bufio.NewScanner(os.Stdin)
			for i := 0; i < len(q); i++ {
				fmt.Print(q[i] + ": ")
				scanner.Scan()

				text := scanner.Text()
				if len(text) != 0 {
					fmt.Println(text)
					paramArr = append(paramArr, text)
				} else {
					paramArr = append(paramArr, "")
				}
			}
			println("Add Options (y/n): ")

			for j := 0; j < len(o); j++ {
				fmt.Print(o[j] + ": ")
				scanner.Scan()

				text := scanner.Text()
				if text == "y" {
					fmt.Println(text)
					paramArr = append(paramArr, o[j])
				} else if text == "n" {
					continue
				}
			}

			bcParameters := "{\"rpcUser\":" + paramArr[0] + ",\"rpcPass\":" + paramArr[1] + ",\"difficulty\":" + paramArr[2] + ",\"options\":[" + strings.Join(paramArr[4:], ",") + "],\"extras\":[" + paramArr[3] + "]}"

			param := "{\"servers\":" + fmt.Sprintf("%s", server) + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"},\"params\":" + bcParameters + "}"
			wsEmitListen(serverAddr, command, param)

		}
		println(checkServ(server))
	},
}

func init() {
	buildCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(buildCmd)
}
