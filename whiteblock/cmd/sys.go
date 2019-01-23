package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type List struct {
	Results []struct {
		Series []struct {
			Columns []string        `json:"columns"`
			Values  [][]interface{} `json:"values"`
		}
	} `json:"results"`
}

var sysCMD = &cobra.Command{
	Use:   "sys <command>",
	Short: "Run SYS commands.",
	Long: "\nSys will allow the user to get information and run SYS commands.\n",
	Run: PartialCommand,
}

var sysTestCMD = &cobra.Command{
	Use:   "test <command>",
	Short: "SYS test commands.",
	Long: "\nSys test will allow the user to get infromation and run SYS tests.\n",
	Run:PartialCommand,
}

var testStartCMD = &cobra.Command{
	Use:   "start <minimum latency> <minimum completion percentage> <number of assets to send> <asset sends per block>",
	Short: "Starts propagation test.",
	Long: `
Sys test start will start the propagation test. It will wait for the signal start time, have nodes send messages at the same time, and require to wait a minimum amount of time then check receivers with a completion rate of minimum completion percentage. 

Format: <wait time> <min complete percent> <number of tx>
Params: Time in seconds, percentage, number of transactions

`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckArguments(args,4,4)
		jsonRpcCallAndPrint("sys::start_test",args)
		return
	},
}

var testResultsCMD = &cobra.Command{
	Use:   "results <test number>",
	Short: "Get results from a previous test.",
	Long: `
Sys test results pulls data from a previous test or tests and outputs as csv.

Format: <test number>
Params: Test number

	`,

	Run: func(cmd *cobra.Command, args []string) {
		CheckArguments(args,1,1)
		results,err := jsonRpcCall("sys::get_recent_test_results",args)
		if err != nil{
			PrintErrorFatal(err)
		}
		result,ok := results.(map[string]interface{})
		if !ok {
			panic(1)
		}
		type List struct {
			Results []struct {
				Series []struct {
					Columns []string        `json:"columns"`
					Values  [][]interface{} `json:"values"`
				}
			} `json:"results"`
		}
		rc := result["results"].([]interface{})[0]
		s := rc.(map[string]interface{})["series"]
		sc := s.([]interface{})[0].(map[string]interface{})
		c := sc["columns"].([]interface{})
		v := sc["values"].([]interface{})

		for i := 0; i < len(v); i ++ {
			fmt.Printf("[%d]\n",i)
			vv := v[i].([]interface{})
			for j := 0; j < len(vv); j++ {
				fmt.Printf("\t%v: %v\n",c[j],vv[j])
			} 
		}

		/*var l List
		json.Unmarshal(results, &l)
		r := l.Results
		rc := r[0]
		s := rc.Series
		sc := s[0]
		c := sc.Columns
		v := sc.Values[0]*/

		/*for i := 0; i < len(c); i++ {
			fmt.Println("\t" + c[i] + ": " + fmt.Sprint(v[i]) + " type is: " + fmt.Sprint(reflect.TypeOf(v[i])))
		}*/
	},
}

var testPinnacleCMD = &cobra.Command{
	Use:   "pinnacle <wait time> <min complete percent>",
	Short: "Run the pinnacle test series.",
	Long: `
Sys test pinnacle will run the propagation test.

Format: <test number>
Params: Test number

	`,

	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		var emptyArgs []string
		var pinArgs []string
		var avgTestTime string

		xS := "200"
		y := 500
		z := 0

		buildCmd.Run(buildCmd, emptyArgs)

		time.Sleep(5 * time.Second)

		wg.Add(1)

		pinArgs = append(pinArgs, args[0])
		pinArgs = append(pinArgs, args[1])
		pinArgs = append(pinArgs, xS)

		for i := 0; i < len(pinArgs); i++ {
			fmt.Println(pinArgs[i])
		}

		// sysStartServerAddr := "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		sysStartServerAddr := serverAddr
		pinnacleServerAddr := serverAddr

		go func() {
			param := "{\"waitTime\":" + pinArgs[0] + ",\"minCompletePercent\":" + pinArgs[1] + ",\"numberOfTransactions\":" + pinArgs[2] + "}"
			fmt.Println(sysStartServerAddr)
			wsEmitListen(sysStartServerAddr, "sys::start_test", param)
			wg.Done()
		}()

		fmt.Println("test start completed")
		fmt.Println("waiting to get results")

		for {
			time.Sleep(10 * time.Second)
			runningTest := wsEmitListen(sysStartServerAddr, "state::is_running", "")
			fmt.Println(runningTest)
			fmt.Println(wsEmitListen(sysStartServerAddr, "state::what_is_running", ""))
			// how long does it take to sync sys::get_recent_test_results. there is an index error because there is no data.
			if runningTest == "false" {
				// fmt.Println("getting test results")
				command := "sys::get_recent_test_results"
				results := []byte(wsEmitListen(sysStartServerAddr, command, args[0]))
				var result map[string]interface{}
				err := json.Unmarshal(results, &result)
				if err != nil {
					panic(err)
				}

				var l List
				json.Unmarshal(results, &l)
				r := l.Results
				rc := r[0]
				s := rc.Series
				sc := s[0]
				c := sc.Columns
				v := sc.Values[0]

				avgTestTime = fmt.Sprint(v[1])

				for j := 0; j < len(v); j++ {
					fmt.Println(string(j) + " -- " + fmt.Sprint(c[j]) + " : " + fmt.Sprint(v[j]))
				}
				fmt.Println("v[9] = " + fmt.Sprint(v[9]))
				break
			}
		}

		fmt.Println(avgTestTime)

		xI, err := strconv.Atoi(xS)
		if err != nil {
			panic(err)
		}

		for {
			fmt.Println("going inside pinnacle loops")

			fmt.Println(pinnacleServerAddr)

			fmt.Println("chaging pinArgs[2] to -> " + xS)
			pinArgs[2] = xS
			for k := 0; k < len(pinArgs); k++ {
				fmt.Println(pinArgs[k])
			}

			fmt.Println("build prev command")
			prevBuild, _ := readPrevCmdFile()
			fmt.Println(prevBuild)
			wsEmitListen(pinnacleServerAddr, "build", prevBuild)
			time.Sleep(5 * time.Second)

			go func() {
				fmt.Println("running test again")
				param := "{\"waitTime\":" + pinArgs[0] + ",\"minCompletePercent\":" + pinArgs[1] + ",\"numberOfTransactions\":" + pinArgs[2] + "}"
				fmt.Println(pinnacleServerAddr)
				wsEmitListen(pinnacleServerAddr, "sys::start_test", param)
			}()

			for {
				time.Sleep(10 * time.Second)
				runningTest := wsEmitListen(sysStartServerAddr, "state::is_running", "")
				fmt.Println(runningTest)
				fmt.Println(wsEmitListen(sysStartServerAddr, "state::what_is_running", ""))
				// how long does it take to sync sys::get_recent_test_results. there is an index error because there is no data.
				if runningTest == "false" {
					fmt.Println("getting test results")
					command := "sys::get_recent_test_results"
					results := []byte(wsEmitListen(sysStartServerAddr, command, args[0]))
					var result map[string]interface{}
					err := json.Unmarshal(results, &result)
					if err != nil {
						panic(err)
					}

					var l List
					json.Unmarshal(results, &l)
					r := l.Results
					rc := r[0]
					s := rc.Series
					sc := s[0]
					c := sc.Columns
					v := sc.Values[0]

					for j := 0; j < len(v); j++ {
						fmt.Println(string(j) + " -- " + fmt.Sprint(c[j]) + " : " + fmt.Sprint(v[j]))
					}
					fmt.Println("v[9] = " + fmt.Sprint(v[9]))

					//unit test this

					fmt.Println(fmt.Sprint(v[1]) + " =?= " + fmt.Sprint(avgTestTime))
					if fmt.Sprint(v[1]) != avgTestTime {
						avgTestTime = fmt.Sprint(v[1])
						if fmt.Sprint(v[9]) == "0" {
							fmt.Println("0 success code")
							if y > 50 {
								xI = (xI - y)
								xS = fmt.Sprintf("%d", xI)
								z = xI

								fmt.Println(xI)
								fmt.Println(y)
								fmt.Println(xS)
								fmt.Println("getting out of 0 y>50")
								break
							} else {
								fmt.Println(xI)
								fmt.Println(y)
								fmt.Println(xS)
								fmt.Println("getting out of 0 y<50")
								break
							}
						} else if fmt.Sprint(v[9]) == "1" {
							fmt.Println("1 success code")
							xI = (xI + y)
							y = y / 2
							xS = fmt.Sprintf("%d", xI)
							z = xI

							fmt.Println(xI)
							fmt.Println(y)
							fmt.Println(xS)
							fmt.Println("getting out of 1")
							break
						}
					} else {
						fmt.Println("avgTestTime is equal")
						break
					}
				}
			}
			if y <= 50 {
				break
			}
			// for j := 0; j < len(v); j++ {
			// 	fmt.Println(string(j) + " -- " + " : " + fmt.Sprint(v[j]))
			// }

		}
		fmt.Println(z)
	},
}

func init() {
	testStartCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	// testResultsCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	sysTestCMD.AddCommand(testStartCMD, testResultsCMD, testPinnacleCMD)
	sysCMD.AddCommand(sysTestCMD)
	RootCmd.AddCommand(sysCMD)
}
