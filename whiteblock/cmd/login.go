package cmd

import (
	util "github.com/whiteblock/cli/whiteblock/util"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetRawProfileFromJwt(jwt string) ([]byte, error) {
	body := strings.NewReader("")
	req, err := http.NewRequest("GET", util.ApiBaseURL+"/agent", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(buf.String())
	}
	return []byte(buf.String()), nil
}

var loginCmd = &cobra.Command{
	Hidden: true,
	Use:    "login <jwt> [biome id]",
	Short:  "Authorize the cli using jwt ",
	Long:   "\nGives the user the ability to specify a jwt, within a file, to be used for authentication\n Can be given a file path or a jwt\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)

		jwt, err := ioutil.ReadFile(args[0])
		if err != nil {
			jwt = []byte(args[0])
		}
		rawProfile, err := GetRawProfileFromJwt(string(jwt))
		if err != nil {
			util.PrintStringError("Given jwt is invalid")
			util.PrintErrorFatal(err)
		}
		util.WriteStore("jwt", jwt)
		util.WriteStore("profile", rawProfile)

		if len(args) == 2 {
			util.WriteStore("biome", []byte(args[1]))
		}

		LoadProfile()
		err = LoadBiomeAddress()
		if err != nil {
			util.DeleteStore("jwt")
			util.DeleteStore("profile")
			util.DeleteStore("biome")
			util.PrintErrorFatal(err)
		}

		fmt.Println("Login Success")
		fmt.Printf("Connected to endpoint: %s\n", serverAddr)
	},
}

var logoutCmd = &cobra.Command{
	Aliases: []string{"logout"},
	Hidden:  true,
	Use:     "logoff",
	Short:   "Remove all auth stored",
	Long:    "\nDeletes all stored auth\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.DeleteStore("jwt")
		util.DeleteStore("profile")
		util.DeleteStore("biome")
		cmd.Println("You have been logged off successfully")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(logoutCmd)
}
