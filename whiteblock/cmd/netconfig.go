package cmd

import (
	"github.com/spf13/cobra"
)

var ()

var netconfigCmd = &cobra.Command{
	Use:     "netconfig <command>",
	Aliases: []string{"emulate"},
	Short:   "Network conditions",
	Long: `
Netconfig will introduce persisting network conditions for testing.
	
	bandwidth <amount> <bandwidth type>	Specifies the bandwidth of the network [bps|Kbps|Mbps|Gbps];
	delay <amount> 				Specifies the latency to add [ms];
	loss <percent>				Specifies the amount of packet loss to add [%];
	
	`,

	Run: func(cmd *cobra.Command, args []string) {

		// command := "netem"
		// jsonRpcCallAndPrint("netem",args)
	},
}

var netconfigSetCmd = &cobra.Command{
	Use:     "set [flags]",
	Aliases: []string{"config", "configure"},
	Short:   "Set network conditions",
	Long: `
Netconfig will introduce persisting network conditions for testing.
	
	--bandwidth <amount>	Specifies the bandwidth of the network in mbps;
	--delay <amount> 				Specifies the latency to add [ms];
	--loss <percent>				Specifies the amount of packet loss to add [%];
	
	`,

	Run: func(cmd *cobra.Command, args []string) {

		// command := "netem"
		// jsonRpcCallAndPrint("netem",args)
	},
}

var netconfigClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"off"},
	Short:   "Turn off network conditions",
	Long: `

	
	`,

	Run: func(cmd *cobra.Command, args []string) {

		// jsonRpcCallAndPrint("netem", args)
	},
}

func init() {
	// latencyCmd.Flags().BoolVarP(&randomPing, "random", "r", false, "apply random latency")

	netconfigCmd.AddCommand(netconfigSetCmd, netconfigClearCmd)

	RootCmd.AddCommand(netconfigCmd)
}
