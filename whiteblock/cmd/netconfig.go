package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	util "../util"
)

var (
	limitFlag int
	lossFlag  float64
	delayFlag int
	rateFlag  int
)

/*type NetConfig struct {
	Servers []int
	NetInfo map[string]interface{}
}*/

var netconfigCmd = &cobra.Command{
	Use:     "netconfig <command>",
	Aliases: []string{"emulate"},
	Short:   "Network conditions",
	Long: `
Netconfig will introduce persisting network conditions for testing.
	
	bandwidth <amount> <bandwidth type>	Specifies the bandwidth of the network [bps|kbps|mbps|gbps];
	delay <amount> 				Specifies the latency to add [ms];
	loss <percent>				Specifies the amount of packet loss to add [%%];
	
	`,

	Run: util.PartialCommand,
}

var netconfigSetCmd = &cobra.Command{
	Use:     "set <node> [flags]",
	Aliases: []string{"config", "configure"},
	Short:   "Set network conditions",
	Long: `
Netconfig set will introduce persisting network conditions for testing to a specific node. Please indicate the proper flags with the amount to set.
	
	--bandwidth <amount>	Specifies the bandwidth of the network in mbps;
	--delay <amount> 				Specifies the latency to add [ms];
	--loss <percent>				Specifies the amount of packet loss to add [%%];
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(args, 1, 1)
		previousBuild,err := getPreviousBuild()
		if err != nil{
			util.PrintErrorFatal(err)
		}
		serverID := previousBuild.Servers[0]


		netInfo := make(map[string]interface{})
		node, err := strconv.Atoi(args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}

		netInfo["node"] = node
		if limitFlag != 1000 {
			netInfo["limit"] = limitFlag
		}
		if lossFlag > 0.0 {
			netInfo["loss"] = lossFlag
		}
		if delayFlag > 0 {
			netInfo["delay"] = delayFlag
		}
		if rateFlag > 0 {
			rate := strconv.Itoa(rateFlag)
			rate = rate + "mbps"
			netInfo["rate"] = rate
		}
		networkConf := []interface{}{
			serverID,
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
	
	--bandwidth <amount>	Specifies the bandwidth of the network in mbps;
	--delay <amount> 				Specifies the latency to add [ms];
	--loss <percent>				Specifies the amount of packet loss to add [%%];
	
	`,

	Run: func(cmd *cobra.Command, args []string) {
		netInfo := make(map[string]interface{})
		previousBuild,err := getPreviousBuild()
		if err != nil{
			util.PrintErrorFatal(err)
		}
		serverID := previousBuild.Servers[0]
		if err != nil {
			fmt.Println("conversion error, invalid type for server")
			return
		}

		if limitFlag != 1000 {
			netInfo["limit"] = limitFlag
		}
		if lossFlag > 0.0 {
			netInfo["loss"] = lossFlag
		}
		if delayFlag > 0 {
			netInfo["delay"] = (delayFlag*1000) / 2
		}
		if rateFlag > 0 {
			rate := strconv.Itoa(rateFlag)
			rate = rate + "mbps"
			netInfo["rate"] = rate
		}

		networkConf := []interface{}{
			serverID,
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
		previousBuild,err := getPreviousBuild()
		if err != nil{
			util.PrintErrorFatal(err)
		}
		serverID := previousBuild.Servers[0]
		jsonRpcCallAndPrint("netem_delete", []interface{}{serverID})
	},
}

func init() {
	netconfigSetCmd.Flags().IntVarP(&limitFlag, "limit", "m", 1000, "sets packet limit")
	netconfigSetCmd.Flags().Float64VarP(&lossFlag, "loss", "l", 0.0, "sets packet loss")
	netconfigSetCmd.Flags().IntVarP(&delayFlag, "delay", "d", 0, "sets latency")
	netconfigSetCmd.Flags().IntVarP(&rateFlag, "bandwidth", "b", 0, "sets the bandwidth")

	netconfigAllCmd.Flags().IntVarP(&limitFlag, "limit", "m", 1000, "sets packet limit")
	netconfigAllCmd.Flags().Float64VarP(&lossFlag, "loss", "l", 0.0, "sets packet loss")
	netconfigAllCmd.Flags().IntVarP(&delayFlag, "delay", "d", 0, "sets latency")
	netconfigAllCmd.Flags().IntVarP(&rateFlag, "bandwidth", "b", 0, "sets the bandwidth")

	netconfigCmd.AddCommand(netconfigSetCmd, netconfigAllCmd, netconfigClearCmd)

	RootCmd.AddCommand(netconfigCmd)
}
