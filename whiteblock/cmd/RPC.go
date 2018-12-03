package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc <command>",
	Short: "Rpc interacts with the blockchain",
	Long: `
RPC allows the user to add their own RPC calls and run those RPC commands.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := args[0]
		param := strings.Join(args[1:], " ")
		wsEmitListen(serverAddr, command, param)
	},
}

// this will have the user add their rpc config file to the directory ~/RPC/
var addrpcCmd = &cobra.Command{
	Use:   "add <file name>",
	Short: "Rpc add will add the rpc calls to the backend",
	Long: `
RPC add will add the rpc calls to the backend which will allow those rpc calls to be made to the blockchain. You should be in the directory of where the file is located.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			out, err := exec.Command("bash", "-c", "./whiteblock get data block -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		args = append(args, usr.HomeDir+"/Downloads/")

		cp := "cp " + args[0] + "/" + args[1] + " ~/RPC/"

		out, err := exec.Command("bash", "-c", cp).Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
	},
}

func init() {
	rpcCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	rpcCmd.AddCommand(addrpcCmd)

	RootCmd.AddCommand(rpcCmd)
}
