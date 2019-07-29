package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
	"golang.org/x/crypto/ssh"
	"io"
	"syscall"
	"unsafe"
)

var getLogCmd = &cobra.Command{
	Use:     "log <node>",
	Aliases: []string{"logs"},
	Short:   "Log will dump data pertaining to the node.",
	Long: `
Get stdout and stderr from a node.

Params: node number

Response: stdout and stderr of the blockchain process
	`,
	//tail -f --zero-terminated /output.log
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		testNetId := build.GetPreviousBuildID()
		n := util.CheckAndConvertInt(args[0], "node number")

		if !util.GetBoolFlagValue(cmd, "follow") {
			util.JsonRpcCallAndPrint("log", map[string]interface{}{
				"testnetId": testNetId,
				"node":      n,
				"lines":     util.GetIntFlagValue(cmd, "tail"),
			})
			return
		}
		//Forward the output from tail -f
		nodes := GetNodes()
		util.CheckIntegerBounds(cmd, "node", n, 0, len(nodes)-1)

		client, err := util.NewSshClient(fmt.Sprintf(nodes[n].IP))
		if err != nil {
			util.PrintErrorFatal(err)
		}
		defer client.Close()

		session, err := client.GetSession()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		defer session.Close() //Open up a session

		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		ws := &winsize{}

		retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
			uintptr(syscall.Stdin),
			uintptr(syscall.TIOCGWINSZ),
			uintptr(unsafe.Pointer(ws)))

		if int(retCode) == -1 {
			panic(errno)
		}

		if err := session.RequestPty("xterm", int(ws.Row), int(ws.Col), modes); err != nil {
			util.PrintErrorFatal(err)
		}
		var outReader io.Reader
		outReader, err = session.StdoutPipe()
		if err != nil {
			util.PrintErrorFatal(err)
		}

		err = session.Start("tail -f --zero-terminated /output.log")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		scanner := bufio.NewScanner(outReader)
		scanner.Split(bufio.ScanLines)

		for {
			scanner.Scan()
			util.Print(scanner.Text())
		}
	},
}

var getLogAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all of the logs",
	Long: `Gets all of the logs
	`,
	Run: func(cmd *cobra.Command, args []string) {
		nodes := GetNodes()
		for i := range nodes {
			util.JsonRpcCallAndPrint("log", map[string]interface{}{
				"testnetId": build.GetPreviousBuildID(),
				"node":      i,
				"lines":     util.GetIntFlagValue(cmd, "tail"),
			})
		}
	},
}

func init() {
	getLogCmd.Flags().IntP("tail", "t", -1, "Get only the last x lines")
	getLogCmd.Flags().BoolP("follow", "f", false, "output appended data as the file grows")
	getCmd.AddCommand(getLogCmd)

	getLogAllCmd.Flags().IntP("tail", "t", -1, "Get only the last x lines")
	getLogCmd.AddCommand(getLogAllCmd)
}
