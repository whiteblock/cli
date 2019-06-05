package cmd

import (
	util "../util"
	"fmt"
	"github.com/spf13/cobra"
)

/*EXAMPLE OF QUERYING MATTS API

func GetRawProfileFromJwt(jwt string) ([]byte, error) {
	body := strings.NewReader("")
	req, err := http.NewRequest("GET", util.ApiBaseURL+"/agent", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth, err := util.CreateAuthNHeader()//get the jwt
	if err != nil {
		log.Println(err)
	} else {
		req.Header.Set("Authorization", auth) //If there is an error, dont send this header for now
	}
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
*/

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "[Description]",
	Long: `
	[In depth description]
`,

	Run: func(cmd *cobra.Command, args []string) {
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		//CODE GOES HERE
		fmt.Println(testnetID) //Remove this line
	},
}

func init() {
	//UNCOMMENT TO ADD THE COMMAND
	//RootCmd.AddCommand(sqlCmd)
}
