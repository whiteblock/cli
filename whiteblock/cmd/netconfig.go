package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var netropyCmd = &cobra.Command{
	Use:     "netconfig <engine number> <path number> <command>",
	Aliases: []string{"emulate"},
	Short:   "Network conditions",
	Long: `
Netconfig will introduce persisting network conditions for testing.
	
	latency 			Specifies the latency to add [ms];
	packetloss 			Specifies the amount of packet loss to add [%];
	bandwidth			Specifies the bandwidth of the network [bps|Kbps|Mbps|Gbps];
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"

		if len(args) < 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock netconfig -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		msg := "engine " + args[0] + " path " + args[1] + " " + strings.Join(args[2:], " ")

		wsEmitListen(serverAddr, command, msg)
	},
}

var emulationOnCmd = &cobra.Command{
	Use:   "on <engine number>",
	Short: "Turn on emulation",

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"

		if len(args) < 1 || len(args) > 1 {
			out, err := exec.Command("bash", "-c", "./whiteblock netconfig -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		msg := "engine " + args[0] + " on"

		wsEmitListen(serverAddr, command, msg)
	},
}

var emulationOffCmd = &cobra.Command{
	Use:   "off <engine number>",
	Short: "Turn off emulation",

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"

		if len(args) < 1 || len(args) > 1 {
			out, err := exec.Command("bash", "-c", "./whiteblock netconfig -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		msg := "engine " + args[0] + " off"

		wsEmitListen(serverAddr, command, msg)
	},
}

var latencyCmd = &cobra.Command{
	Use:     "latency <engine number> <path number> <amount>",
	Aliases: []string{"lat"},
	Short:   "Set latency",
	Long: `
Latency will introduce delay to the network. You will specify the amount of latency in ms.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"

		if len(args) < 3 || len(args) > 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock netconfig -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		msg1 := "engine " + args[0] + " path " + args[1] + " set delay constant " + args[2] + " port 1 to port 2"
		msg2 := "engine " + args[0] + " path " + args[1] + " set delay constant " + args[2] + " port 2 to port 1"

		wsEmitListen(serverAddr, command, msg1)
		wsEmitListen(serverAddr, command, msg2)
	},
}

var packetLossCmd = &cobra.Command{
	Use:     "packetloss <engine number> <path number> <percent>",
	Aliases: []string{"loss"},
	Short:   "Set packetloss",
	Long: `
Packetloss will drop packets in the network. You will specify the amount of packet loss in %.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"

		if len(args) < 3 || len(args) > 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock netconfig -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		msg1 := "engine " + args[0] + " path " + args[1] + " set loss random " + args[2] + " port 1 to port 2"
		msg2 := "engine " + args[0] + " path " + args[1] + " set loss random " + args[2] + " port 2 to port 1"

		wsEmitListen(serverAddr, command, msg1)
		wsEmitListen(serverAddr, command, msg2)
	},
}

var bandwCmd = &cobra.Command{
	Use:     "bandwidth <engine number> <path number> <amount> <bandwidth type>",
	Aliases: []string{"bw"},
	Short:   "Set bandwidth",
	Long: `
Bandwidth will constrict the network to the specified bandwidth. You will specify the amount of bandwdth and the type.

Fomat: 
	bandwidth type: bps, Kbps, Mbps, Gbps
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "netconfig"

		if len(args) < 4 || len(args) > 4 {
			out, err := exec.Command("bash", "-c", "./whiteblock netconfig -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		msg1 := "engine " + args[0] + " path " + args[1] + " set bw fixed " + args[2] + args[3]
		msg2 := "engine " + args[2] + " path " + args[3] + " set bw fixed " + args[2] + args[3]

		wsEmitListen(serverAddr, command, msg1)
		wsEmitListen(serverAddr, command, msg2)
	},
}

func init() {
	netropyCmd.AddCommand(emulationOnCmd, emulationOffCmd, latencyCmd, packetLossCmd, bandwCmd)

	RootCmd.AddCommand(netropyCmd)
}
