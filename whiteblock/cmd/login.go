package cmd

import (
	util "../util"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetRawProfileFromJwt(jwt string) ([]byte, error) {
	body := strings.NewReader("")
	req, err := http.NewRequest("GET", "https://api.whiteblock.io/agent", body)
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
	Use:    "login <jwt> [organization] [biome]",
	Short:  "Authorize the cli using jwt ",
	Long:   "\nGives the user the ability to specify a jwt, within a file, to be used for authentication\n Can be given a file path or a jwt\n",
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 3)

		jwt, err := ioutil.ReadFile(args[0])
		if err != nil {
			jwt = []byte(args[0])
		}
		rawOrgKey, err := GetRawProfileFromJwt(string(jwt))
		if err != nil {
			util.PrintStringError("Given jwt is invalid")
			util.PrintErrorFatal(err)
		}
		util.WriteStore("jwt", jwt)
		util.WriteStore("org_key", rawOrgKey)
		switch len(args) {
		case 3:
			util.WriteStore("biome", []byte(args[2]))
			fallthrough
		case 2:
			util.WriteStore("organization", []byte(args[1]))
		}
		LoadOrganizationApiKey()
		err = LoadBiomeAddress()
		if err != nil {
			util.DeleteStore("jwt")
			util.DeleteStore("org_key")
			util.DeleteStore("profile")
			util.DeleteStore("biome")
			util.DeleteStore("organization")
			util.PrintErrorFatal(err)
		}

		fmt.Println("Login Success")
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
		util.DeleteStore("organization")
		cmd.Println("You have been logged off successfully")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(logoutCmd)
}
