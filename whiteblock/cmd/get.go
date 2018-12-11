package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func cwFile(path, data string) {
	time := time.Now().UTC().String()
	time = strings.Replace(time, " ", "", -1)

	file, err := os.Create(path + "/dataset_" + time + ".txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

var getCmd = &cobra.Command{
	Use:   "get <command>",
	Short: "Get server and network information.",
	Long: `
Get will ouput server and network information and statstics.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock get -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var getServerCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"servers"},
	Short:   "Get server information.",
	Long: `
Server will ouput server information.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_servers"

		println(wsEmitListen(serverAddr, command, ""))
	},
}

var getNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Nodes will show all nodes in the network.",
	Long: `
Nodes will output all of the nodes in the current network.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_nodes"
		println(wsEmitListen(serverAddr, command, ""))
	},
}

var getDataCmd = &cobra.Command{
	Use:   "data <command>",
	Short: "Data will pull data from the network and output into a file.",
	Long: `
Data will pull specific or all block data from the network and output into a file. You will specify the directory where the file will be downloaded.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock get data -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var dataByTimeCmd = &cobra.Command{
	Use:   "time <start time> <end time> [path]",
	Short: "Data time will pull data from the network and output into a file.",
	Long: `
Data time will pull block data from the network from a given start and end time and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || len(args) > 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock get data time -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		} else if len(args) == 2 {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			args = append(args, usr.HomeDir+"/Downloads/")
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":" + args[0] + ",\"endTime\":" + args[1] + ",\"startBlock\":0,\"endBlock\":0}"
		data := wsEmitListen(serverAddr, command, param)

		cwFile(args[2], data)
	},
}

var dataByBlockCmd = &cobra.Command{
	Use:   "block <start block> <end block> [path]",
	Short: "Data block will pull data from the network and output into a file.",
	Long: `
Data block will pull block data from the network from a given start and end block and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || len(args) > 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock get data block -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		} else if len(args) == 2 {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			args = append(args, usr.HomeDir+"/Downloads/")
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":0,\"endTime\":0,\"startBlock\":" + args[0] + ",\"endBlock\":" + args[1] + "}"
		data := wsEmitListen(serverAddr, command, param)

		cwFile(args[2], data)
	},
}

var dataAllCmd = &cobra.Command{
	Use:   "all [path]",
	Short: "All will pull data from the network and output into a file.",
	Long: `
Data all will pull all data from the network and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			out, err := exec.Command("bash", "-c", "./whiteblock get data all -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		} else if len(args) == 0 {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			args = append(args, usr.HomeDir+"/Downloads/")
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "all_stats"
		data := wsEmitListen(serverAddr, command, "")

		cwFile(args[0], data)
	},
}

var getStatsCmd = &cobra.Command{
	Use:   "stats <command>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock get stats -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var statsByTimeCmd = &cobra.Command{
	Use:   "time <start time> <end time>",
	Short: "Get stastics by time",
	Long: `
Stats time will allow the user to get statistics by specifying a start time and stop time (unix time stamp).

Params: Unix time stamps
Format: <start unix time stamp> <end unix time stamp>

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			out, err := exec.Command("bash", "-c", "./whiteblock get stats time -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":" + args[0] + ",\"endTime\":" + args[1] + ",\"startBlock\":0,\"endBlock\":0}"
		data := wsEmitListen(serverAddr, command, param)
		println(data)
	},
}

var statsByBlockCmd = &cobra.Command{
	Use:   "block <start block> <end block>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats block will allow the user to get statistics regarding the network.

Params: Block numbers
Format: <start block number> <end block number>

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			out, err := exec.Command("bash", "-c", "./whiteblock get stats block -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "stats"
		param := "{\"startTime\":0,\"endTime\":0,\"startBlock\":" + args[0] + ",\"endBlock\":" + args[1] + "}"
		data := wsEmitListen(serverAddr, command, param)
		println(data)
	},
}

var statsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all stastics of a blockchain",
	Long: `
Stats all will allow the user to get all the statistics regarding the network.

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "all_stats"
		data := wsEmitListen(serverAddr, command, "")
		println(data)
	},
}

func init() {
	getCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	getServerCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	getNodesCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getDataCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	dataByTimeCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	dataByBlockCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	dataAllCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getStatsCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsByTimeCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsByBlockCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsAllCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getCmd.AddCommand(getServerCmd, getNodesCmd, getDataCmd, getStatsCmd)
	getDataCmd.AddCommand(dataByTimeCmd, dataByBlockCmd, dataAllCmd)
	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsAllCmd)

	RootCmd.AddCommand(getCmd)
}
