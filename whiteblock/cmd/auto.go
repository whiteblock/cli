package cmd

import (
	"encoding/json"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"sort"
	"strconv"
	"time"
)

var autoCmd = &cobra.Command{
	Aliases: []string{},
	Use:     "auto <node> <command> [params]",
	Short:   "send queries",
	Long: `Automatically send json_rpc queries to a node in the background. <command> is the name of the json rpc call to be made. 
	You can use +account,+tx_hash,+number,+hex,+block_hash,+block_number as magic string parameters to be filled in with randomized appropriate values.
	+tx_hash random tx hash; only works after you call wb tx start stream
	+account random account
	+number random base 10 number
	+hex random base16 number
	+block_hash random block hash
	+block_number random block number
	Examples:
	object parameter with an array:
	wb auto 0 eth_sendTransaction -i 1000000 '{"from":"+account","to":"+account","gas":"0x76c0","gasPrice":"0x9184e72a000","value":"+hex","data":"0x00"}'
	simple parameters within an array:
	wb auto -i 100000 0 eth_getBalance  +account latest
`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, -1)
		node, err := strconv.Atoi(args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		interval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		sampleSize, err := cmd.Flags().GetInt("sample-size")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		errorChecking, err := cmd.Flags().GetBool("full-error-checking")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		maxNumErrMsgs, err := cmd.Flags().GetUint("max-num-err-msgs")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		recordErrMsgs, err := cmd.Flags().GetBool("disable-error-recording")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		params := []interface{}{}
		if len(args) > 2 {
			for _, arg := range args[2:] {
				var param interface{}
				err = json.Unmarshal([]byte(arg), &param)
				if err != nil {
					param = arg //if it is not json, then it is a string
				}
				params = append(params, param)
			}
		}

		util.JsonRpcCallAndPrint("setup_load", []interface{}{map[string]interface{}{
			"node": node,
			"name": fmt.Sprintf("node%d:%s", node, args[1]),
			"settings": map[string]interface{}{
				"targetDelay":   interval,
				"sampleSize":    sampleSize,
				"maxNumErrMsgs": maxNumErrMsgs,
				"recordErrMsgs": !recordErrMsgs,
			},
			"call":       args[1],
			"arguments":  params,
			"errorCheck": errorChecking,
		}})
	},
}

var autoKillCmd = &cobra.Command{
	Use:     "kill",
	Aliases: []string{"stop"},
	Short:   "stop an auto routine",
	Long: `
Kill an auto routine.
`,
	Run: func(cmd *cobra.Command, args []string) {
		forced, err := cmd.Flags().GetBool("force")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if forced {
			util.CheckArguments(cmd, args, 1, 1)
			util.JsonRpcCallAndPrint("state::force_stop_sub_routine", []interface{}{args[0]})
		} else {
			util.JsonRpcCallAndPrint("state::kill_sub_routines", args)
		}
	},
}

var autoCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean a stoped auto routine",
	Long: `
clean a stoped auto routine
`,
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("state::clean_sub_routines", args)
	},
}

var autoPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "purge all of the running autos",
	Long: `
Gracefully stops and removes all of the currently running auto routines.
Most users do not need to call this as it happens automatically on the next build.
`,
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("state::purge_all_sub_routines", args)
	},
}

var getAutoCmd = &cobra.Command{
	Use:     "auto",
	Aliases: []string{"routines"},
	Short:   "Check auto QPS",
	Long:    "Get the QPS of the currently running automated queries",
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("state::sub_routines", []string{})
	},
}

var getAutoErrorsCmd = &cobra.Command{
	Use:     "errors",
	Aliases: []string{"error"},
	Short:   "Check auto errors",
	Long:    "Get the most recent error messages",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			util.JsonRpcCallAndPrint("state::all_sub_routines_errors", []string{})
			return
		}
		util.JsonRpcCallAndPrint("state::sub_routine_errors", args)

	},
}

func createAutoGraph() ([]ui.Drawable, error) {
	res, err := util.JsonRpcCall("state::sub_routines_stats", []string{})
	if err != nil {
		return nil, err
	}
	if len(res.(map[string]interface{})) == 0 {

		return nil, fmt.Errorf("nothing to show")
	}
	log.WithFields(log.Fields{"length": len(res.(map[string]interface{}))}).Trace("fetched the stats")
	width, _ := getTermSize()
	y := 0

	increment := 10

	/**sorting mechanics**/
	keys := []string{}
	for routine, _ := range res.(map[string]interface{}) {
		keys = append(keys, routine)
	}
	tmpKeys := sort.StringSlice(keys)
	tmpKeys.Sort()
	sortedKeys := []string(tmpKeys)

	objects := []ui.Drawable{}
	for _, routine := range sortedKeys {
		data := res.(map[string]interface{})[routine].(map[string]interface{})
		//render the data
		var points []map[string]float64
		stats := data["stats"].(map[string]interface{})
		historical := stats["historical"]
		tmp, err := json.Marshal(historical)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(tmp, &points)
		if err != nil {
			return nil, err
		}
		/*Create the graph*/
		plot := widgets.NewPlot()
		plot.Title = routine
		plot.LineColors[0] = ui.ColorRed
		plot.LineColors[1] = ui.ColorYellow
		plot.LineColors[2] = ui.ColorGreen
		//plot.LineColors = []ui.Color{,ui.ColorYellow,ui.ColorRed}
		plot.AxesColor = ui.ColorWhite
		plot.DataLabels = []string{"eps", "qps", "sps"}
		//plot.LineColors[0] = ui.ColorGreen
		//plot.HorizontalScale = (50/len(points)) + 1
		plot.Data = make([][]float64, 3)
		for _, sample := range points {
			secs := sample["time_microseconds"] / 1000000.0
			eps := sample["errors"] / secs
			qps := sample["queries"] / secs
			sps := sample["successes"] / secs
			plot.Data[0] = append(plot.Data[0], eps)
			plot.Data[1] = append(plot.Data[1], qps)
			plot.Data[2] = append(plot.Data[2], sps)
		}
		//plot.Data = [][]float64{[]float64{1,2},[]float64{3,4},[]float64{5,6}}
		plot.DrawDirection = widgets.DrawLeft

		plot.SetRect(0, y, int(width*2), increment+y)

		objects = append(objects, plot)

		/*Create the overall table*/
		table := widgets.NewTable()
		table.Rows = [][]string{
			[]string{"successes", fmt.Sprintf("%v", int64(data["successes"].(float64)))},
			[]string{"errors", fmt.Sprintf("%v", int64(data["errors"].(float64)))},
			[]string{"success rate", fmt.Sprintf("%v", data["successRate"])},
		}
		table.TextStyle = ui.NewStyle(ui.ColorWhite)
		table.RowSeparator = true
		table.BorderStyle = ui.NewStyle(ui.ColorWhite)
		table.SetRect(int(width*2)+1, y, int(width*2)+1+int(width), 7+y)
		table.FillRow = true
		/*table3.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
		table3.RowStyles[2] = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
		table3.RowStyles[3] = ui.NewStyle(ui.ColorYellow)*/
		objects = append(objects, table)
		y += increment
	}
	return objects, nil
}

var getAutoDetailedCmd = &cobra.Command{
	Use:     "detail",
	Aliases: []string{"details", "detailed"},
	Short:   "Check the progress of auto queries in detail",
	Long:    "Check the progress of auto queries in detail",
	Run: func(cmd *cobra.Command, args []string) {
		graphIt, err := cmd.Flags().GetBool("graph")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if !graphIt {
			util.JsonRpcCallAndPrint("state::sub_routines_stats", []string{})
			return
		}
		if err = ui.Init(); err != nil {
			util.PrintErrorFatal(err)
		}
		defer ui.Close()

		go func() {
			for {

				plots, err := createAutoGraph()
				if err != nil {
					ui.Close()
					util.PrintErrorFatal(err)
				}
				ui.Render(plots...)
				time.Sleep(time.Second)
			}

		}()

		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}

	},
}

func init() {
	autoCmd.Flags().Bool("full-error-checking", false, "Check for errors other than just connectivity errors (default false)")
	autoCmd.Flags().IntP("interval", "i", 50000, "Send interval in microseconds")
	autoCmd.Flags().IntP("sample-size", "s", 200, "auto stats sample size")
	autoCmd.Flags().Uint("max-num-err-msgs", 5, "the max history of error messages to keep (default 5)")
	autoCmd.Flags().Bool("disable-error-recording", false, "disable the recording of error messages."+
		" (default false) Saves ~1us per error")
	autoKillCmd.Flags().BoolP("force", "f", false, "force kill/stop the routine (this may cause a crash)")

	getAutoDetailedCmd.Flags().Bool("graph", false, "show an interactive graph of the results")
	autoCmd.AddCommand(autoKillCmd, autoCleanCmd, autoPurgeCmd)

	getAutoCmd.AddCommand(getAutoDetailedCmd, getAutoErrorsCmd)
	getCmd.AddCommand(getAutoCmd)
	RootCmd.AddCommand(autoCmd)

}
