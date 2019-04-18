package util

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func PartialCommand(cmd *cobra.Command, args []string) {
	fmt.Println("\nNo command given. Please choose a command from the list below.")
	cmd.Help()
	return
}

func CheckAndConvertInt(num string, name string) int {
	out, err := strconv.ParseInt(num, 0, 32)
	if err != nil {
		InvalidInteger(name, num, true)
	}
	return int(out)
}

func CheckAndConvertInt64(num string, name string) int64 {
	out, err := strconv.ParseInt(num, 0, 64)
	if err != nil {
		InvalidInteger(name, num, true)
	}
	return out
}

/*
   Write writes data to a file, creating it if it doesn't exist,
   deleting and recreating it if it does.
*/
func Write(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0664)
}

/*
	Sends an http request and returns the body. Gives an error if the http request failed
	or returned a non success code.
*/
func HttpRequest(method string, url string, bodyData string) ([]byte, error) {
	//log.Println("URL IS "+url)
	body := strings.NewReader(bodyData)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	return buf.Bytes(), nil
}

func UnrollStringSliceToMapStringString(slices []string, delim string) (map[string]string, error) {
	out := map[string]string{}
	for _, slice := range slices {
		pair := strings.SplitN(slice, delim, 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf(`Missing "%s" delimiter in flag value.`, delim)
		}
		out[pair[0]] = pair[1]
	}
	return out, nil
}

func UnrollStringSliceToMapIntString(slices []string, delim string) (map[int]string, []string, error) {
	out := map[int]string{}
	noDelimRes := []string{}
	for _, slice := range slices {
		pair := strings.SplitN(slice, delim, 2)
		if len(pair) != 2 {
			noDelimRes = append(noDelimRes, slice)
			continue
		}
		key, err := strconv.Atoi(pair[0])
		if err != nil {
			return nil, nil, err
		}

		out[key] = pair[1]
	}
	return out, noDelimRes, nil
}
