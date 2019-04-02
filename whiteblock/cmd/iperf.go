package cmd

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"os"
	"io"
	"bufio"
	"syscall"
	"unsafe"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	util "../util"
)

type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

func getWidth() uint {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno)
    }
    return uint(ws.Col)
}

var (
	bw          string
	testTime    string
	udpEnabled  bool
	dualEnabled bool
)

func CaptureAndOutput(r io.Reader) {
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanLines)
    for scanner.Scan() {
        m := scanner.Text()
        fmt.Println(m)
    }	
}

func PadString(str string,target int) string {
	out := str
	for i := len(str); i  < target ; i++ {
		out += " "
	}
	return out
}
//Display the contents of two readers, line by line, together
func CaptureAndDisplayTogether(r1 io.Reader,r2 io.Reader,offset int,label1 string,label2 string) {
	minWidth := 160

	scanner1 := bufio.NewScanner(r1)
	scanner1.Split(bufio.ScanLines)

	scanner2 := bufio.NewScanner(r2)
	scanner2.Split(bufio.ScanLines)

	var red1 bool 
	var red2 bool
	var txt1 string
	var txt2 string

	
	width := int(getWidth())
	centerSize := int(width/20)
	panelSize := int(width/2) - centerSize
	padding := PadString("",centerSize)

	counter := 0
	if width > minWidth {
		fmt.Printf("%s%s%s\n",PadString(label1,panelSize),padding,PadString(label2,panelSize))
	}
	for{
		

		red1 = scanner1.Scan()
		if counter >= offset {
			red2 = scanner2.Scan()
		}
		
		if !(red1 || red2){
			break
		}

		if red1 {
			txt1 = scanner1.Text()
		}else{
			txt1 = ""
		}

		if red2 && counter >= offset {
			txt2 = scanner2.Text()
		}else{
			txt2 = ""
		}
		txt1 = PadString(txt1,panelSize)
		txt2 = PadString(txt2,panelSize)
		if width > minWidth{
			fmt.Printf("%s%s%s\n",txt1,padding,txt2)
		}else{
			fmt.Printf("%s:%s\n%s:%s\n",label1,txt1,label2,txt2)
		}
		
		counter++
	}
}

var iPerfCmd = &cobra.Command{
	Use:   "iperf <sending node> <receiving node>",
	Short: "iperf will show network conditions.",
	Long: `

Iperf will show the user network conditions and other data. This command will establish the sending node as a server and the receiving node as a client node. They will send packets and at the end of the test, the output will give bandwidth, transfer size, and other relevant

Format: <sending node> <receiving node>
Params: sending node, receiving node
	`,

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		util.CheckArguments(cmd,args,2,2)
		spinner := Spinner{}
		spinner.SetText("Setting Up Iperf")
		spinner.Run(100)

		nodes,err := GetNodes()
		if err != nil{
			util.PrintErrorFatal(err)
		}

		sendingNodeNumber, err := strconv.Atoi(args[0])
		if err != nil {
			util.InvalidArgument(args[0])
			cmd.Help()
			return
		}
		receivingNodeNumber, err := strconv.Atoi(args[1])
		if err != nil {
			util.InvalidArgument(args[1])
			cmd.Help()
			return
		}
		if sendingNodeNumber >= len(nodes) {
			util.PrintStringError("Sending node number too high")
			os.Exit(1)
		}

		if receivingNodeNumber >= len(nodes) {
			util.PrintStringError("Receiving node number too high")
			os.Exit(1)
		}

		var outReader1 		io.Reader
		var outReader2 		io.Reader
		var awaitReaders 	sync.WaitGroup 
		awaitReaders.Add(2)
		wg.Add(2)
		// command to run iperf as a server
		go func() {
			defer wg.Done()

			iPerfcmd := "iperf3 -s "
			if udpEnabled {
				iPerfcmd = iPerfcmd + "-u "
			}

			iPerfcmd = iPerfcmd + fmt.Sprintf(nodes[sendingNodeNumber].IP) + " -1"

			client, err := util.NewSshClient(fmt.Sprintf(nodes[sendingNodeNumber].IP))
			if err != nil {
				util.PrintErrorFatal(err)
			}
			defer client.Close()

			client.Run("pkill -9 iperf3")//Kill iperf if it is running

			session,err := client.GetSession()
			if err != nil {
				util.PrintErrorFatal(err)
			}
			defer session.Close()//Open up a session

			modes := ssh.TerminalModes{
			    ssh.ECHO:          0,
			    ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			    ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}

			if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
			    util.PrintErrorFatal(err)
			}

			outReader1,err = session.StdoutPipe()
			if err != nil {
				util.PrintErrorFatal(err)
			}
			awaitReaders.Done()
			//go CaptureAndOutput(outReader)

			err = session.Start(iPerfcmd)
			if err != nil {
				util.PrintErrorFatal(err)
			}
			session.Wait()

		}()

		go func() {
			// command to run iperf as a client
			time.Sleep(500 * time.Millisecond)
			defer wg.Done()

			iPerfcmd := "iperf3 -c "
			if udpEnabled {
				iPerfcmd = iPerfcmd + " -u "
			}
			if bw != "" && udpEnabled {
				_, err := strconv.Atoi(bw)
				if err != nil {
					fmt.Println("Invalid format given for bandwidth flag.")
					return
				}
				iPerfcmd = iPerfcmd + " -b " + bw
			} else if bw != "" && !udpEnabled {
				fmt.Println("udp needs to be enabled to set bandwidth.")
			}
			if dualEnabled {
				iPerfcmd = iPerfcmd + " -d "
			}

			iPerfcmd = iPerfcmd + fmt.Sprintf(nodes[sendingNodeNumber].IP)

			client, err := util.NewSshClient(fmt.Sprintf(nodes[receivingNodeNumber].IP))
			if err != nil {
				util.PrintErrorFatal(err)
			}
			defer client.Close()

			client.Run("pkill -9 iperf3")
			spinner.Kill()
			session,err := client.GetSession()
			if err != nil {
				util.PrintErrorFatal(err)
			}
			defer session.Close()//Open up a session

			modes := ssh.TerminalModes{
			    ssh.ECHO:          0,
			    ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			    ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}

			if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
			    util.PrintErrorFatal(err)
			}

			outReader2,err = session.StdoutPipe()
			if err != nil {
				util.PrintErrorFatal(err)
			}
			awaitReaders.Done()
			awaitReaders.Wait()
			//go CaptureAndOutput(outReader)

			err = session.Start(iPerfcmd)
			if err != nil {
				util.PrintErrorFatal(err)
			}
			session.Wait()
		}()
		awaitReaders.Wait()
		go CaptureAndDisplayTogether(outReader1,outReader2,3,"SERVER","CLIENT")
		wg.Wait()
	},
}

func init() {
	iPerfCmd.Flags().StringVarP(&bw, "bandwidth", "b", "", "set target bandwidth in bits/sec (default 1 Mbit/sec); requires udp enabled")
	iPerfCmd.Flags().BoolVarP(&dualEnabled, "dualtest", "d", false, "enable bidirectional test simultaneously")
	iPerfCmd.Flags().StringVarP(&testTime, "time", "t", "", "how long to run test for")
	iPerfCmd.Flags().BoolVarP(&udpEnabled, "udp", "u", false, "enable udp")

	RootCmd.AddCommand(iPerfCmd)
}
