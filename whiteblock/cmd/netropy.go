package cmd

import (
	"github.com/spf13/cobra"
)

var (
	latency    int
	packetloss float32
)

var netropyCmd = &cobra.Command{
	Use:   "netropy",
	Short: "Network conditions",
	Long: `Netropy will introduce persisting network conditions for testing.
	
	latency 			Specifies the latency to add [ms];
	packetloss 			Specifies the amount of packet loss to add [%];
	`,

	Run: func(cmd *cobra.Command, args []string) {
		// add curl command

	},
}

func init() {
	netropyCmd.Flags().IntVarP(&latency, "latency", "l", 10, "latency")
	netropyCmd.Flags().Float32VarP(&packetloss, "packetloss", "p", 0.001, "packetloss")

	RootCmd.AddCommand(netropyCmd)
}
