package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	dir     string
	cwd     []string
	lastcmd string
)

var sshCmd = &cobra.Command{
	Use:   "ssh <server> <node>",
	Short: "SSH into an existing container.",
	Long: `
SSH will allow the user to go into the contianer where the specified node exists.

Response: stdout of the command
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		if len(args) != 2 {
			println("Invalid number of arguments given")
			out, err := exec.Command("bash", "-c", "./whiteblock ssh -h").Output()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		dir = "Server" + args[0] + "-Node" + args[1] + ":"
		cwd = append(cwd, "/")
		cwd = append(cwd, "root/")

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGSTOP)
		go func() {
			for {
				sig := <-sigs
				fmt.Println(sig)
				command := "exec"

				param1 := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + "ps aux | grep \\\"" + lastcmd + "\\\" |grep -v grep| awk '{print $2}' | sort | tail -n 1\"}"
				println(param1)
				pid := wsEmitListen(serverAddr, command, param1)
				pid = strings.Replace(pid, "\n", "", -1)

				param2 := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + "kill " + pid + "\"}"
				resp := wsEmitListen(serverAddr, command, param2)
				println(resp, lastcmd, ": process has been terminated")

				fmt.Print("\r\n"+dir+strings.Join(cwd[:], ""), "$ ")
			}
		}()

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("SSH Server: " + args[0] + " Node: " + args[1])
		fmt.Println("---------------------------------------------------------------")

		for {
			//print cwd path
			fmt.Print(dir+strings.Join(cwd[:], ""), "$ ")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			textarg := strings.Split(text, " ")

			//to quit
			if text == "q" {
				os.Exit(1)
			}

			//cd handling
			if textarg[0] == "cd" {
				if len(textarg) == 1 {
					cwd = nil
					cwd = append(cwd, "/")
				} else if textarg[1] == ".." {
					if len(cwd) > 0 {
						cwd = append(cwd[:len(cwd)-1])
					}
				} else if textarg[1] != ".." {
					command := "exec"
					param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + "[ -d " + strings.Join(cwd, "") + textarg[1] + "/ ] && echo $? \"}"
					outerr := wsEmitListen(serverAddr, command, param)

					if outerr[0] == '0' {
						cwd = append(cwd, textarg[1]+"/")
					} else {
						println("directory does not exist")
					}
				}
			} else if len(textarg[0]) == 0 {
				continue
			} else if textarg[0] == "yes" || textarg[0] == "kill" {
				println("command not found")
			} else {
				// println(cwd)
				lastcmd = textarg[0]

				command := "exec"
				param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + "bash -c \\\"cd " + strings.Join(cwd, "") + " && " + text + "\\\"\"}"
				println(param)
				out := wsEmitListen(serverAddr, command, param)
				println(out)
			}
		}
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}

/*
fixes/bugs:
- have to account for "cd ../" or "cd ..{anything}"
	- this will append to the cwd improperly

	*- fix this by error handling by checking if the directory actually exists before appending
*/
