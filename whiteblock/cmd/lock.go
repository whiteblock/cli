package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"time"
)

var acquireCmd = &cobra.Command{
	Use:     "acquire",
	Aliases: []string{"lock"},
	Short:   "Wait until a lock has been acquired on the endpoint",
	Long: `
This call will block until a unique lock has been acquired on the endpoint
	`,
	Run: func(cmd *cobra.Command, args []string) {
		once, err := cmd.Flags().GetBool("once")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		conf.RPCRetries = 1
		if once {
			util.JsonRpcCallAndPrint("lock", []interface{}{})
			return
		}
		for {
			res, err := util.JsonRpcCall("lock", []interface{}{})
			if err == nil {
				fmt.Println(util.Prettyp(res.(string)))
				break
			}
			time.Sleep(time.Second * 10)
		}

	},
}

var forceUnlockCmd = &cobra.Command{
	Use:   "force-unlock",
	Short: "Forces an unlock",
	Long:  "\nForces an unlock, might be dangerous, but is useful if it is stuck in lock mode\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("unlock", []interface{}{})
	},
}

func init() {
	acquireCmd.Flags().Bool("once", false, "try to acquire the lock only once")
	RootCmd.AddCommand(acquireCmd, forceUnlockCmd)

}
