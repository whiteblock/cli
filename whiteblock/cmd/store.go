package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"sync"
)

var storeCommand = &cobra.Command{
	Hidden: true,
	Use:    "store",
	Short:  "store",
	Long:   `store`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		wg := sync.WaitGroup{}
		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		var extras map[string]interface{}
		err = json.Unmarshal(data, &extras)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		for key, val := range extras {
			wg.Add(1)
			go func(key string, val interface{}) {
				defer wg.Done()
				_, err := util.JsonRpcCall("set_extra", []interface{}{key, val})
				if err != nil {
					util.PrintErrorFatal(err)
				}
			}(key, val)
		}
		wg.Wait()
		util.Print("done")
	},
}

func init() {
	RootCmd.AddCommand(storeCommand)
}
