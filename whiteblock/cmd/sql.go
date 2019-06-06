package cmd

import (
	"bytes"
	"encoding/json"
	"log"

	"fmt"
	"io/ioutil"
	"net/http"

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
	Kind string
	Etag string
	Tables []table

}

type table struct {
	Kind string
	Id string
	TableReference struct {
		ProjectId string
		DatasetId string
		TableId string
	}
}

type metricsResponse struct {
	Schema schema `json:"schema"`
	JobReference struct {
		JobID string `json:"jobId"`
	} `json:"jobReference"`
	TotalRows int `json:"totalRows"`
	PageToken string `json:"pageToken"`
	Rows [][]interface{} `json:"rows"`
	Error errr `json:"error"`

}

type schema struct {
}

type errr struct {
}


var sqlCmd = &cobra.Command{
	Use:   "sql <command>",
	Short: "",
	Long: `
    `,
	Run: util.PartialCommand,
}

var sqlTableListCmd = &cobra.Command{
	Use:   "list",
	Short: "Gets a list of current tables in the database",
	Long: `
sql list will return a list of current tables in the database	

Response: JSON representation of the table list in the database
	`,

	Run: func(cmd *cobra.Command, args []string) {
		payload := []byte{}

		id, err := cmd.Flags().GetInt("organization-id")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		data, err := apiRequest(fmt.Sprintf("/organizations/%d/dw/tables", id), "GET", payload)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		var response interface{}

		err = json.Unmarshal(data, &response)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		fmt.Println("    TABLES     ")

		log.Println(prettypi(response))
	},
}

type SqlQueryRequestPayload struct {
	Q string `json:"q"`
}

var sqlQueryCmd = &cobra.Command{
	Use:   "query <query>",
	Short: "Runs SQL command to retrieve structured log data",
	Long: `
This command will run a SQL query to the database to retrieve structured log data

Format: whiteblock sql query <SQL query>
	`,

	Run: func(cmd *cobra.Command, args []string) {
		query := SqlQueryRequestPayload{Q: args[0]}
		payload, err := json.Marshal(query)
		if err != nil {
			util.PrintErrorFatal(err)
		}


		id, err := cmd.Flags().GetInt("organization-id")
		if err != nil {
			util.PrintErrorFatal(err)
		}

		data, err := apiRequest(fmt.Sprintf("/organizations/%d/dw/metrics", id), "POST", payload)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		var metrics metricsResponse

		fmt.Println(data)

		err = json.Unmarshal(data, &metrics)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		fmt.Println("    METRICS     ")

		log.Println(prettypi(metrics))
	},
}

func init() {
	sqlTableListCmd.Flags().IntP("organization-id", "i", 10, "api request returns the specified organization's data")
	sqlQueryCmd.Flags().IntP("organization-id", "i", 10, "api request returns the specified organization's data")
	sqlCmd.AddCommand(sqlTableListCmd, sqlQueryCmd)
	RootCmd.AddCommand(sqlCmd)
}

func apiRequest(path string, method string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", util.ApiBaseURL, path), bytes.NewReader(body))
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
