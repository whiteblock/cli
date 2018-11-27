package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	gc "github.com/rthornton128/goncurses"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh <server> <node> <command> ",
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

		var mutex = &sync.Mutex{}
		mutex.Lock()

		scr, err := gc.Init()
		if err != nil {
			log.Fatal("init:", err)
		}
		defer gc.End()

		gc.Echo(false)

		scr.Println("Type characters to have them appear on the screen.")
		scr.Println("Press 'q' to exit.")
		scr.Println()

		// Accept input concurrently via a goroutine and connect a channel
		in := make(chan gc.Char)
		ready := make(chan bool)
		go func(w *gc.Window, ch chan<- gc.Char) {
			for {
				// Block until all write operations are complete
				<-ready
				// Send typed character down the channel (which is blocking
				// in the main loop)
				ch <- gc.Char(w.GetChar())
			}
		}(scr, in)

		// Once a character has been received on the 'in' channel the
		// 'ready' channel will block until it receives another piece of data.
		// This happens only once the received character has been written to
		// the screen. The 'in' channel then blocks on the next loop until
		// another 'true' is sent down the 'ready' channel signaling to the
		// input goroutine that it's okay to receive input
		for {
			var c gc.Char
			select {
			case c = <-in: // blocks while waiting for input from goroutine
				scr.Print(string(c))
				scr.Refresh()
			case ready <- true: // sends once above block completes
			}
			// Exit when 'q' is pressed
			if c == gc.Char('q') {
				mutex.Unlock()
				break
			}
		}

		command := "exec"
		param := "{\"server\":" + args[0] + ",\"node\":" + args[1] + ",\"command\":\"" + strings.Join(args[2:], " ") + "\"}"

		wsEmitListen(serverAddr, command, param)
	},
}

func init() {
	sshCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(sshCmd)
}
