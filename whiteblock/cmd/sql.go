package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	util "github.com/whiteblock/cli/whiteblock/util"
)

// Query from the userdata API
/*
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

var sqlTableListCmd = &cobra.Command{
	Use:   "sql list",
	Short: "Gets a list of current tables in the database",
	Long: `
sql list will return a list of current tables in the database	

Response: JSON representation of the table list in the database
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

var sqlQueryCmd = &cobra.Command{
	Use:   "sql <query>",
	Short: "Runs SQL command to retrieve structured log data",
	Long: `
This command will run a SQL query to the database to retrieve structured log data
	
Format: whiteblock sql <SQL query>
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
	RootCmd.AddCommand(sqlTableListCmd)
	RootCmd.AddCommand(sqlQueryCmd)
}
