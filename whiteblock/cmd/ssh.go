package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	dir string
)

func getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

var sshCmd = &cobra.Command{
	Use:   "ssh <server> <node> //<command> ",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		if len(args) != 3 {
			println("Invalid number of arguments given")
			out, err := exec.Command("bash", "-c", "./whiteblock ssh -h").Output()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		// command := "exec"
		// param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + strings.Join(args[2:], " ") + "\"}"

		// wsEmitListen(serverAddr, command, param)

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Simple Shell")
		fmt.Println("---------------------------------------------------------------")
		on := true

		for on {
			fmt.Print(getDir() + " -> ")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)

			command := "exec"
			param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + text + "\"}"

			wsEmitListen(serverAddr, command, param)

			if text == "q" {
				break
			}

			fmt.Println(text)
		}

	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
