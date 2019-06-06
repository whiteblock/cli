package cmd

import (
	"bytes"
	"encoding/json"
	"log"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/whiteblock/cli/whiteblock/util"

	"github.com/spf13/cobra"
)

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

		var response struct {
			Kind   string `json:"kind"`
			Etag   string `json:"etag"`
			Tables []struct {
				Kind           string `json:"kind"`
				Id             string `json:"id"`
				TableReference struct {
					ProjectId string `json:"projectId"`
					DatasetId string `json:"datasetId"`
					TableId   string `json:"tableId"`
				} `json:"tableReference"`
			} `json:"tables"`
		}

		err = json.Unmarshal(data, &response)
		if err != nil {
			util.PrintErrorFatal(err)
		}

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

		var metrics struct {
			Schema       interface{} `json:"schema"`
			JobReference struct {
				JobID string `json:"jobId"`
			} `json:"jobReference"`
			TotalRows int             `json:"totalRows"`
			PageToken string          `json:"pageToken"`
			Rows      [][]interface{} `json:"rows"`
			Error     interface{}     `json:"error"`
		}

		fmt.Println(data)

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
	sqlCmd.AddCommand(sqlTableListCmd, sqlQueryCmd)
	RootCmd.AddCommand(sqlCmd)
}

func apiRequest(path string, method string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", util.ApiBaseURL, path), bytes.NewReader(body))
	if err != nil {
		util.PrintErrorFatal(err)
	}

	auth, err := util.CreateAuthNHeader() //get the jwt
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
