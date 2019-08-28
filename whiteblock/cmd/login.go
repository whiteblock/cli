package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetProfileFromJwt(jwt string) (Profile, error) {
	var out Profile
	body := strings.NewReader("")
	req, err := http.NewRequest("GET", conf.APIURL+"/agent", body)
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return out, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out, fmt.Errorf(buf.String())
	}

	return out, json.Unmarshal([]byte(buf.String()), &out)
}

var loginCmd = &cobra.Command{
	Use:   "login <jwt> [biome id]",
	Short: "Authorize the cli using jwt ",
	Long:  "\nGives the user the ability to specify a jwt, within a file, to be used for authentication\n Can be given a file path or a jwt\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)

		jwt, err := ioutil.ReadFile(args[0])
		if err != nil {
			jwt = []byte(args[0])
		}
		prof, err := GetProfileFromJwt(string(jwt))
		if err != nil {
			util.PrintStringError("Given jwt is invalid")
			util.PrintErrorFatal(err)
		}
		util.Set("jwt", string(jwt))
		util.Set("profile", prof)

		if len(args) == 2 {
			util.Set("biome", args[1])
		}

		LoadProfile()
		err = LoadBiomeAddress()
		if err != nil {
			util.Delete("jwt")
			util.Delete("profile")
			util.Delete("biome")
			util.PrintErrorFatal(err)
		}

		util.Print("Login Success")
		fmt.Printf("Connected to endpoint: %s\n", conf.ServerAddr)
	},
}

var logoutCmd = &cobra.Command{
	Aliases: []string{"logout"},
	Use:     "logoff",
	Short:   "Remove all auth stored",
	Long:    "\nDeletes all stored auth\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.Delete("jwt")
		util.Delete("profile")
		util.Delete("biome")
		util.Print("You have been logged off successfully")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(logoutCmd)
}
