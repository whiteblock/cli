package cmd

import (
	"encoding/json"
	"log"

	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	util "github.com/whiteblock/cli/whiteblock/util"

	"github.com/spf13/cobra"
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

type tableResponse struct {
	kind string
	etag string
	tables []table

}

type table struct {
	kind string
	id string
	tableReference tableReference
}

type tableReference struct {
	projectId string
	datasetId string
	tableId string
}

type metricsResponse struct {
	schema schema
	jobReference jobReference
	totalRows int
	pageToken string
	rows rows
	error errr

}

type schema struct {
}

type jobReference struct {
	jobId string
}

type rows struct {
}

type errr struct {
}




var sqlTableListCmd = &cobra.Command{
	Use:   "sql list",
	Short: "Gets a list of current tables in the database",
	Long: `
sql list will return a list of current tables in the database	

Response: JSON representation of the table list in the database
	`,

	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("organization-id")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		data, err := apiRequest(fmt.Sprintf("/organizations/%d/dw/tables", id), "GET")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		var response interface{}

		err = json.Unmarshal(data, &response)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		log.Println(prettypi(response))
	},
}

var sqlQueryCmd = &cobra.Command{
	Use:   "sql query <query>",
	Short: "Runs SQL command to retrieve structured log data",
	Long: `
This command will run a SQL query to the database to retrieve structured log data

Format: whiteblock sql query <SQL query>
	`,

	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("organization-id")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		data, err := apiRequest(fmt.Sprintf("/organizations/%d/dw/metrics", id), "POST")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		var metrics metricsResponse

		err = json.Unmarshal(data, &metrics)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		log.Println(prettypi(metrics))
	},
}

func init() {
	sqlTableListCmd.Flags().IntP("organization-id", "i", 10, "api request returns the specified organization's data")
	sqlQueryCmd.Flags().IntP("organization-id", "i", 10, "api request returns the specified organization's data")
	RootCmd.AddCommand(sqlTableListCmd)
	RootCmd.AddCommand(sqlQueryCmd)
}

func apiRequest(path string, method string) ([]byte, error) {

	body := strings.NewReader("")
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", util.ApiBaseURL, path), body)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	auth, err := util.CreateAuthNHeader()//get the jwt
	if err != nil {
		util.PrintErrorFatal(err)
	} else {
		request.Header.Set("Authorization", auth) //If there is an error, dont send this header for now
	}
	request.Close = true

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	defer resp.Body.Close()


	return ioutil.ReadAll(resp.Body)
}
