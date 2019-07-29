package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
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

		delay := util.GetIntFlagValue(cmd, "delay")
		rate := util.GetIntFlagValue(cmd, "bandwidth")
		limit := util.GetIntFlagValue(cmd, "limit")

		util.CheckArguments(cmd, args, 1, 1)

		netInfo := make(map[string]interface{})
		node := util.CheckAndConvertInt(args[0], "node")

		netInfo["node"] = node
		if limit != 1000 {
			netInfo["limit"] = limit
		}
		if lossFlag > 0.0 {
			netInfo["loss"] = lossFlag
		}
		if delay > 0 {
			netInfo["delay"] = delay * 1000
		}
		if rate > 0 {
			rateStr := strconv.Itoa(rate)
			rateStr = rateStr + "mbps"
			netInfo["rate"] = rateStr
		}
		networkConf := []interface{}{
			build.GetPreviousBuildID(),
			netInfo,
		}

		util.JsonRpcCallAndPrint("netem", networkConf)
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

		delay := util.GetIntFlagValue(cmd, "delay")
		rateFlag := util.GetIntFlagValue(cmd, "bandwidth")
		limit := util.GetIntFlagValue(cmd, "limit")
		if limit != 1000 {
			netInfo["limit"] = limit
		}
		if lossFlag > 0.0 {
			netInfo["loss"] = lossFlag
		}
		if delay > 0 {
			netInfo["delay"] = (delay * 1000) / 2
		}
		if rateFlag > 0 {
			rate := strconv.Itoa(rateFlag)
			rate = rate + "mbps"
			netInfo["rate"] = rate
		}

		networkConf := []interface{}{
			build.GetPreviousBuildID(),
			netInfo,
		}

		util.JsonRpcCallAndPrint("netem_all", networkConf)
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
		util.JsonRpcCallAndPrint("netem_delete", []interface{}{build.GetPreviousBuildID()})
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
		util.JsonRpcCallAndPrint("netem_get", []interface{}{build.GetPreviousBuildID()})
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
		outArgs := []interface{}{build.GetPreviousBuildID()}
		if len(args) == 1 {
			outArgs = append(outArgs, args[0])
		}
		util.JsonRpcCallAndPrint("get_outages", outArgs)
	},
}

var netconfigGetPartitionsCmd = &cobra.Command{
	Use: "partitions",
	//Aliases: []string{"blocked", "disconnected"},
	Short: "Get the network partitions",
	Long:  "\nGets the current network partitions\n",

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		util.JsonRpcCallAndPrint("get_partitions", []interface{}{build.GetPreviousBuildID()})
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
		util.JsonRpcCallAndPrint("remove_outage", []interface{}{
			build.GetPreviousBuildID(),
			util.CheckAndConvertInt(args[0], "node1"),
			util.CheckAndConvertInt(args[1], "node2")})
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
		util.JsonRpcCallAndPrint("make_outage", []interface{}{
			build.GetPreviousBuildID(),
			util.CheckAndConvertInt(args[0], "node1"),
			util.CheckAndConvertInt(args[1], "node2")})
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
		nodes := []int{}
		for i, arg := range args {
			nodes = append(nodes, util.CheckAndConvertInt(arg, fmt.Sprintf("argument %d", i)))
		}
		util.JsonRpcCallAndPrint("partition_outage", []interface{}{build.GetPreviousBuildID(), nodes})
	},
}

var netconfigMarryCmd = &cobra.Command{
	Use:   "marry",
	Short: "Remove any outages",
	Long: `
Remove any outages and allow connections between all nodes
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		util.JsonRpcCallAndPrint("remove_all_outages", []interface{}{build.GetPreviousBuildID()})
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
