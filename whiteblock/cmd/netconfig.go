package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	latency    int
	packetloss float32
)

var netropyCmd = &cobra.Command{
	Use:     "netconfig <engine number> <path number> <command>",
	Aliases: []string{"emulate"},
	Short:   "Network conditions",
	Long: `Netconfig will introduce persisting network conditions for testing.
	
	latency 			Specifies the latency to add [ms];
	packetloss 			Specifies the amount of packet loss to add [%];
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"
		msg := "engine " + args[0] + " path " + args[1] + " " + strings.Join(args[2:], " ")

		wsSendCmd(serverAddr, command, msg)
	},
}

func init() {
	netropyCmd.Flags().IntVarP(&latency, "latency", "l", 10, "latency")
	netropyCmd.Flags().Float32VarP(&packetloss, "packetloss", "p", 0.001, "packetloss")

	RootCmd.AddCommand(netropyCmd)
}
