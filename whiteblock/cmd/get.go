package cmd

import (
	"fmt"
	"log"
	"os"
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
		println("\nNo command given. Please choose a command from the list above.\n")
		cmd.Help()
		return
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

		println(prettyp(wsEmitListen(serverAddr, command, "")))
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

var getRunningCmd = &cobra.Command{
	Use:   "running",
	Short: "Running will check if a test is running.",
	Long: `
Running will check whether or not there is a test running and get the name of the currently running test.

Response: true or false, on whether or not a test is running; The name of the test or nothing if there is not a test running.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "state::is_running"
		command2 := "state::what_is_running"
		println(wsEmitListen(serverAddr, command1, ""))
		println(wsEmitListen(serverAddr, command2, ""))

	},
}

var getLogCmd = &cobra.Command{
	Use:   "log <node number>",
	Short: "Log will dump data pertaining to the node.",
	Long: `
Get stdout and stderr from a node.

Params: node number

Response: stdout and stderr of the blockchain process
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "log"
		param := "{\"server\":" + fmt.Sprintf("%d", server) + ",\"node\":" + args[0] + "}"
		println(wsEmitListen(serverAddr, command, param))
	},
}

var getNetworkDefaultsCmd = &cobra.Command{
	Use:   "default <blockchain>",
	Short: "Default gets the blockchain params.",
	Long: `
Get the blockchain specific parameters for a deployed blockchain.

Params: <blockchain>
Format: The blockchain to get the build params of

Response: The params as a list of key value params, of name and type respectively
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "get_defaults"
		println(wsEmitListen(serverAddr, command, args[0]))
	},
}

// var getDataCmd = &cobra.Command{
// 	Use:   "data <command>",
// 	Short: "Data will pull data from the network and output into a file.",
// 	Long: `
// Data will pull specific or all block data from the network and output into a file. You will specify the directory where the file will be downloaded.

// 	`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		println("\nNo command given. Please choose a command from the list above.\n")
// 		cmd.Help()
// 		return
// 	},
// }

// var dataByTimeCmd = &cobra.Command{
// 	Use:   "time <start time> <end time> [path]",
// 	Short: "Data time will pull data from the network and output into a file.",
// 	Long: `
// Data time will pull block data from the network from a given start and end time and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

// Params: Unix time stamps
// Format: <start unix time stamp> <end unix time stamp>

// 	`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		if len(args) < 2 || len(args) > 3 {
// 			println("\nError: Invalid number of arguments given\n")
// 			cmd.Help()
// 			return
// 		} else if len(args) == 2 {
// 			usr, err := user.Current()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			args = append(args, usr.HomeDir+"/Downloads/")
// 		}

// 		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
// 		command := "stats"
// 		param := "{\"startTime\":" + args[0] + ",\"endTime\":" + args[1] + ",\"startBlock\":0,\"endBlock\":0}"
// 		data := wsEmitListen(serverAddr, command, param)

// 		cwFile(args[2], data)
// 	},
// }

// var dataByBlockCmd = &cobra.Command{
// 	Use:   "block <start block> <end block> [path]",
// 	Short: "Data block will pull data from the network and output into a file.",
// 	Long: `
// Data block will pull block data from the network from a given start and end block and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

// Params: Unix time stamps
// Format: <start unix time stamp> <end unix time stamp>

// 	`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		if len(args) < 2 || len(args) > 3 {
// 			println("\nError: Invalid number of arguments given\n")
// 			cmd.Help()
// 			return
// 		} else if len(args) == 2 {
// 			usr, err := user.Current()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			args = append(args, usr.HomeDir+"/Downloads/")
// 		}

// 		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
// 		command := "stats"
// 		param := "{\"startTime\":0,\"endTime\":0,\"startBlock\":" + args[0] + ",\"endBlock\":" + args[1] + "}"
// 		data := wsEmitListen(serverAddr, command, param)

// 		cwFile(args[2], data)
// 	},
// }

// var dataAllCmd = &cobra.Command{
// 	Use:   "all [path]",
// 	Short: "All will pull data from the network and output into a file.",
// 	Long: `
// Data all will pull all data from the network and output into a file. The directory where the file will be downloaded will need to be specified. If no directory is provided, default directory is set to ~/Downloads.

// 	`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		if len(args) > 1 {
// 			println("\nError: Invalid number of arguments given\n")
// 			cmd.Help()
// 			return
// 		} else if len(args) == 0 {
// 			usr, err := user.Current()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			args = append(args, usr.HomeDir+"/Downloads/")
// 		}

// 		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
// 		command := "all_stats"
// 		data := wsEmitListen(serverAddr, command, "")

// 		cwFile(args[0], data)
// 	},
// }

var getStatsCmd = &cobra.Command{
	Use:   "stats <command>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics
	`,

	Run: func(cmd *cobra.Command, args []string) {
		println("\nError: Invalid number of arguments given\n")
		cmd.Help()
		return
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
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
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
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
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

	// getDataCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	// dataByTimeCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	// dataByBlockCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	// dataAllCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getStatsCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsByTimeCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsByBlockCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	statsAllCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd, getNetworkDefaultsCmd, getRunningCmd, getLogCmd)
	// getDataCmd.AddCommand(dataByTimeCmd, dataByBlockCmd, dataAllCmd)
	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsAllCmd)

	RootCmd.AddCommand(getCmd)
}
