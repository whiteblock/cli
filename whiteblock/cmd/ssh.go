package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"io/ioutil"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

type Node []struct {
	ID        int    `json:"id"`
	TestNetID int    `json:"testNetId"`
	Server    int    `json:"server"`
	LocalID   int    `json:"localId"`
	IP        string `json:"ip"`
	Label     string `json:"label"`
}

var sshCmd = &cobra.Command{
	Use:   "ssh <node>",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.
	`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("\nError: Invalid number of arguments given")
			cmd.Help()
			return
		}

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command1 := "nodes"
		out1 := []byte(wsEmitListen(serverAddr, command1, ""))
		var node Node
		json.Unmarshal(out1, &node)
		nodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		val,_ := os.LookupEnv("HOME")
		dat, err := ioutil.ReadFile(val+"/.ssh/id_rsa.pub")

		if err != nil {
			fmt.Println("Run ssh-keygen first!")
			os.Exit(1)
		}
		command2 := "exec"
		param := "{\"server\":" + server + ",\"node\":" + args[0] + ",\"command\":\"service ssh start\"}"
		wsEmitListen(serverAddr, command2, param)

		param = "{\"server\":" + server + ",\"node\":" + args[0] + 
				",\"command\":\"bash -c \\\"echo \\\\\\\""+strings.Trim(string(dat),"\n\t\r\v")+"\\\\\\\">> /root/.ssh/authorized_keys\\\"\"}"
		wsEmitListen(serverAddr, command2, param)

		_, err = exec.Command("bash", "-c", "rm $HOME/.ssh/known_hosts").Output()
		if err != nil {
			fmt.Println("No known hosts")
		}

		err = unix.Exec("/usr/bin/ssh", []string{"ssh", "-o", "StrictHostKeyChecking no", "root@" + fmt.Sprintf(node[nodeNumber].IP)}, os.Environ())
		log.Fatal(err)
		println(nodeNumber)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
