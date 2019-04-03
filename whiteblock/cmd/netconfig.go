package cmd

import (
	util "../util"
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
		testnetId, err := getPreviousBuildId()
		if err != nil {
			util.PrintStringError("No previous build found")
			os.Exit(1)
		}
		jsonRpcCallAndPrint("netem_delete", []interface{}{testnetId})
	},
}

func init() {
	netconfigSetCmd.Flags().IntVarP(&limitFlag, "limit", "m", 1000, "sets packet limit")
	netconfigSetCmd.Flags().Float64VarP(&lossFlag, "loss", "l", 0.0, "Specifies the amount of packet loss to add [%%];")
	netconfigSetCmd.Flags().IntVarP(&delayFlag, "delay", "d", 0, "Specifies the latency to add [ms];")
	netconfigSetCmd.Flags().IntVarP(&rateFlag, "bandwidth", "b", 0, "Specifies the bandwidth of the network in mbps;")

	netconfigAllCmd.Flags().IntVarP(&limitFlag, "limit", "m", 1000, "sets packet limit")
	netconfigAllCmd.Flags().Float64VarP(&lossFlag, "loss", "l", 0.0, "Specifies the amount of packet loss to add [%%];")
	netconfigAllCmd.Flags().IntVarP(&delayFlag, "delay", "d", 0, "Specifies the latency to add [ms];")
	netconfigAllCmd.Flags().IntVarP(&rateFlag, "bandwidth", "b", 0, "Specifies the bandwidth of the network in mbps;")

	netconfigCmd.AddCommand(netconfigSetCmd, netconfigAllCmd, netconfigClearCmd)

	RootCmd.AddCommand(netconfigCmd)
}
