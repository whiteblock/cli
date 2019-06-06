package cmd

import (
	util "github.com/whiteblock/cli/whiteblock/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	limitFlag int
	lossFlag  float64
	delayFlag int
	rateFlag  int
)

var netconfigCmd = &cobra.Command{
	Use:     "netconfig <command>",
	Aliases: []string{"emulate"},
	Short:   "Network conditions",
	Long: `
Netconfig will introduce persisting network conditions for testing.
`,
	Run: util.PartialCommand,
}

var netconfigSetCmd = &cobra.Command{
	Use:     "set <node> [flags]",
	Aliases: []string{"config", "configure"},
	Short:   "Set network conditions",
	Long: `
Netconfig set will introduce persisting network conditions for testing to a specific node. Please indicate the proper flags with the amount to set.
`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}

		netInfo := make(map[string]interface{})
		node, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidInteger("node", args[0], true)
		}

		netInfo["node"] = node
		if limitFlag != 1000 {
			netInfo["limit"] = limitFlag
		}
		if lossFlag > 0.0 {
			netInfo["loss"] = lossFlag
		}
		if delayFlag > 0 {
			netInfo["delay"] = delayFlag * 1000
		}
		if rateFlag > 0 {
			rate := strconv.Itoa(rateFlag)
			rate = rate + "mbps"
			netInfo["rate"] = rate
		}
		networkConf := []interface{}{
			testnetId,
			netInfo,
		}

		jsonRpcCallAndPrint("netem", networkConf)
	},
}

var netconfigAllCmd = &cobra.Command{
	Use:     "all [flags]",
	Aliases: []string{"config", "configure"},
	Short:   "Set network conditions",
	Long: `
Netconfig all will introduce persisting network conditions for testing to all nodes. Please indicate the proper flags with the amount to set.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		netInfo := make(map[string]interface{})
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}

		if limitFlag != 1000 {
			netInfo["limit"] = limitFlag
		}
		if lossFlag > 0.0 {
			netInfo["loss"] = lossFlag
		}
		if delayFlag > 0 {
			netInfo["delay"] = (delayFlag * 1000) / 2
		}
		if rateFlag > 0 {
			rate := strconv.Itoa(rateFlag)
			rate = rate + "mbps"
			netInfo["rate"] = rate
		}

		networkConf := []interface{}{
			testnetId,
			netInfo,
		}

		jsonRpcCallAndPrint("netem_all", networkConf)
	},
}

var netconfigClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"off", "flush", "reset"},
	Short:   "Turn off network conditions",
	Long: `
Netconfig clear will reset all emulation and turn off all persisiting network conditions. 
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		jsonRpcCallAndPrint("netem_delete", []interface{}{testnetId})
	},
}

var netconfigGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"show"},
	Short:   "Get the network conditions",
	Long: `
Netconfig get will fetch the current network conditions
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		jsonRpcCallAndPrint("netem_get", []interface{}{testnetId})
	},
}

var netconfigGetDisconnectsCmd = &cobra.Command{
	Use:     "disconnects [node]",
	Aliases: []string{"blocked", "disconnected"},
	Short:   "Get the blocked connections",
	Long: `
Get a json array of the connections which are blocked. 
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		outArgs := []interface{}{testnetId}
		if len(args) == 1 {
			outArgs = append(outArgs, args[0])
		}
		jsonRpcCallAndPrint("get_outages", outArgs)
	},
}

var netconfigGetPartitionsCmd = &cobra.Command{
	Use: "partitions",
	//Aliases: []string{"blocked", "disconnected"},
	Short: "Get the network partitions",
	Long:  "\nGets the current network partitions\n",

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		/*spinner := Spinner{}
		spinner.SetText("Fetching the network partitions")
		spinner.Run(100)
		defer spinner.Kill()*/
		jsonRpcCallAndPrint("get_partitions", []interface{}{testnetId})
	},
}

var netconfigUncutCmd = &cobra.Command{
	Use:     "uncut <node1> <node2>",
	Aliases: []string{"unblock"},
	Short:   "Allow the given pair of nodes to connect",
	Long: `
Allow the given pair of nodes to connect
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		node1, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidInteger("node1", args[0], true)
		}
		node2, err := strconv.Atoi(args[1])
		if err != nil {
			util.InvalidInteger("node2", args[1], true)
		}
		jsonRpcCallAndPrint("remove_outage", []interface{}{testnetId, node1, node2})
	},
}

var netconfigCutCmd = &cobra.Command{
	Use:     "cut <node1> <node2>",
	Aliases: []string{"block"},
	Short:   "Prevent the given pair of nodes from connecting",
	Long: `
Prevent the given pair of nodes from connecting
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		node1, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidInteger("node1", args[0], true)
		}
		node2, err := strconv.Atoi(args[1])
		if err != nil {
			util.InvalidInteger("node2", args[1], true)
		}
		jsonRpcCallAndPrint("make_outage", []interface{}{testnetId, node1, node2})
	},
}

var netconfigPartitionCmd = &cobra.Command{
	Use: "partition <node1>...",
	//Aliases: []string{"unblock"},
	Short: "Partition the given nodes from the rest of the network",
	Long: `
Partition the given nodes from the rest of the network
	`,

	Run: func(cmd *cobra.Command, args []string) {

		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		nodes := []int{}
		for i, arg := range args {
			node, err := strconv.Atoi(arg)
			if err != nil {
				util.InvalidInteger(fmt.Sprintf("argument %d", i), args[0], true)
			}
			nodes = append(nodes, node)
		}
		/*spinner := Spinner{}
		spinner.SetText("Partition the network")
		spinner.Run(100)
		defer spinner.Kill()*/
		jsonRpcCallAndPrint("partition_outage", []interface{}{testnetId, nodes})
	},
}

var netconfigMarryCmd = &cobra.Command{
	Use: "marry",
	//Aliases: []string{"unblock"},
	Short: "Remove any outages",
	Long: `
Remove any outages and allow connections between all nodes
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		/*spinner := Spinner{}
		spinner.SetText("Putting the network back together")
		spinner.Run(100)
		defer spinner.Kill()*/
		jsonRpcCallAndPrint("remove_all_outages", []interface{}{testnetId})
	},
}

func init() {
	netconfigSetCmd.Flags().IntVarP(&limitFlag, "limit", "m", 1000, "sets packet limit")
	netconfigSetCmd.Flags().Float64VarP(&lossFlag, "loss", "l", 0.0, "Specifies the amount of packet loss to add [%]")
	netconfigSetCmd.Flags().IntVarP(&delayFlag, "delay", "d", 0, "Specifies the latency to add [ms]")
	netconfigSetCmd.Flags().IntVarP(&rateFlag, "bandwidth", "b", 0, "Specifies the bandwidth of the network in mbps")

	netconfigAllCmd.Flags().IntVarP(&limitFlag, "limit", "m", 1000, "sets packet limit")
	netconfigAllCmd.Flags().Float64VarP(&lossFlag, "loss", "l", 0.0, "Specifies the amount of packet loss to add [%]")
	netconfigAllCmd.Flags().IntVarP(&delayFlag, "delay", "d", 0, "Specifies the latency to add [ms]")
	netconfigAllCmd.Flags().IntVarP(&rateFlag, "bandwidth", "b", 0, "Specifies the bandwidth of the network in mbps")

	netconfigGetCmd.AddCommand(netconfigGetDisconnectsCmd, netconfigGetPartitionsCmd)

	netconfigCmd.AddCommand(netconfigSetCmd, netconfigAllCmd, netconfigClearCmd, netconfigGetCmd, netconfigUncutCmd,
		netconfigCutCmd, netconfigPartitionCmd, netconfigMarryCmd)

	RootCmd.AddCommand(netconfigCmd)
}
