package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
)

var sqlCmd = &cobra.Command{
	Use:   "sql <command>",
	Short: "sql runs SQL queries to obtain organization data",
	Long: `
sql runs SQL queries to obtain organization data, specifically metrics and tables
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

		if id == 0 {
			id = getOrgId()
		}

		data, err := apiRequest(fmt.Sprintf("/organizations/%d/dw/tables", id), "GET", payload)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		var tables struct {
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

		err = json.Unmarshal(data, &tables)

		fmt.Println(prettypi(tables))
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

		if id == 0 {
			id = getOrgId()
		}

		data, err := apiRequest(fmt.Sprintf("/organizations/%d/dw/metrics", id), "POST", payload)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		var response metrics
		err = json.Unmarshal(data, &response)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		outRows := make([][]interface{}, 0)
		outRows = append(outRows, response.Rows...)

		for response.PageToken != "" {
			response = response.next(id)
			outRows = append(outRows, response.Rows...)
		}

		fmt.Println(
			prettypi(
				struct {
					Schema interface{}     `json:"schema"`
					Rows   [][]interface{} `json:"rows"`
				}{
					Schema: response.Schema,
					Rows:   outRows,
				},
			),
		)
	},
}

func init() {
	sqlTableListCmd.Flags().IntP("organization-id", "i", 0, "api request returns the specified organization's data")
	sqlQueryCmd.Flags().IntP("organization-id", "i", 0, "api request returns the specified organization's data")
	sqlCmd.AddCommand(sqlTableListCmd, sqlQueryCmd)
	RootCmd.AddCommand(sqlCmd)
}

func apiRequest(path string, method string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", util.GetConfig().APIURL, path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	auth, err := util.CreateAuthNHeader() //get the jwt
	if err != nil {
		return nil, err
	} else {
		request.Header.Set("Authorization", auth) //If there is an error, dont send this header for now
	}
	request.Close = true

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Trace(string(data))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s\nstatus code is %d", string(data), resp.StatusCode)
	}

	return data, nil
}

type metrics struct {
	Schema       interface{} `json:"schema"`
	JobReference struct {
		JobID string `json:"jobId"`
	} `json:"jobReference"`
	TotalRows int             `json:"totalRows"`
	PageToken string          `json:"pageToken"`
	Rows      [][]interface{} `json:"rows"`
	Error     interface{}     `json:"error"`
}

func (m *metrics) next(id int) metrics {
	path := fmt.Sprintf("/organizations/%d/dw/metrics?job_id=%s&page_token=%s", id, m.JobReference.JobID, m.PageToken)
	data, err := apiRequest(path, "GET", []byte{})
	if err != nil {
		util.PrintErrorFatal(err)
	}

	response := metrics{}

	err = json.Unmarshal(data, &response)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	return response
}

func getOrgId() int {
	id, err := apiRequest("/agent", "GET", []byte{})
	if err != nil {
		util.PrintErrorFatal(err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(id, &response)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	return int(response["organization_id"].(float64))
}
