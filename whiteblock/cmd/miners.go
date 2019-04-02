package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	util "../util"
)

var minerCmd = &cobra.Command{
	// Hidden: true,
	Use:   "miner <command>",
	Short: "Run miner commands.",
	Long: `
Send commands pertaining to mining. This will be blockchain specific and will only be supported depending on which blockchain had been built.
	`,
	Run: util.PartialCommand,
}

var minerStartCmd = &cobra.Command{
	Use:   "start [node 1 number] [node 2 number]...",
	Short: "Start Mining",
	Long: `
Send the start mining signal to nodes, may take a while to take effect due to DAG generation. If no arguments are given, all nodes will begin mining.

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to start mining or None for all nodes

Response: The number of nodes which successfully received the signal to start mining`,
	Run: func(cmd *cobra.Command, args []string) {
		spinner := &Spinner{txt:"Starting the miner",die:false}
		spinner.Run(100)
		
		res, err := jsonRpcCall("start_mining", args)
		if err != nil {
			util.PrintErrorFatal(err)
			/*util.PrintStringError("There was an error building the DAG.")
			os.Exit(1)*/
		}
		DagReady := false
		for !DagReady {
			//fmt.Printf("\rDAG is being generated...")
			res, err = jsonRpcCall("get_block_number", []string{})
			if err != nil {
				util.PrintErrorFatal(err)
			}
			blocknum := int(res.(float64))
			if blocknum > 2 {
				DagReady = true
			}
			time.Sleep(time.Millisecond * 50)
		}
		//fmt.Println("\rDAG has been successfully generated.")
		spinner.Kill()
		time.Sleep(time.Millisecond * 100)
		fmt.Println("\rDAG has been successfully generated.")
	},
}

var minerStopCmd = &cobra.Command{
	Use:   "stop [node 1 number] [node 2 number]...",
	Short: "Stop mining",
	Long: `
Send the stop mining signal to nodes

Format: [node 1 number] [node 2 number]...
Params: A list of the nodes to stop mining or None for all nodes


Response: The number of nodes which successfully received the signal to stop mining`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("stop_mining", args)
	},
}

func init() {
	minerCmd.AddCommand(minerStartCmd, minerStopCmd)
	RootCmd.AddCommand(minerCmd)
}
