package util

import (
	"bytes"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestPartialCommand(t *testing.T) {
	//var tests = []struct {
	//	cmd *cobra.Command
	//	args []string
	//}{
	//	{
	//		cmd: &cobra.Command{},
	//		args: []string{},
	//	},
	//	{
	//		cmd: &cobra.Command{},
	//		args: []string{"one"},
	//	},
	//	{
	//		cmd: &cobra.Command{},
	//		args: []string{"one", "two"},
	//	},
	//}

	// TODO How to compare prints?

}

func TestCheckAndConvertInt(t *testing.T) {
	var tests = []struct {
		num      string
		name     string
		expected int
	}{
		{num: "5", name: "test", expected: 5},
		{num: "158348", name: "test", expected: 158348},
		{num: "0", name: "test", expected: 0},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if CheckAndConvertInt(tt.num, tt.name) != tt.expected {
				t.Error("return value of CheckAndConvertInt does not match expected value")
			}
		})
	}
}

func TestCheckAndConvertInt64(t *testing.T) {
	var tests = []struct {
		num      string
		name     string
		expected int64
	}{
		{num: "5", name: "test", expected: 5},
		{num: "158348", name: "test", expected: 158348},
		{num: "0", name: "test", expected: 0},
		{num: "2392347592347", name: "test", expected: 2392347592347},
		{num: "-35", name: "test", expected: -35},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if CheckAndConvertInt64(tt.num, tt.name) != tt.expected {
				t.Error("return value of CheckAndConvertInt does not match expected value")
			}
		})
	}
}

func TestWrite(t *testing.T) {
	path := "/tmp/testWrite"
	data := []byte("blah")

	if Write(path, data) != nil {
		t.Error("return value of Write does not match expected value")
	}
}

func TestHttpRequest(t *testing.T) {
	method := "GET"
	url := "https://this-page-intentionally-left-blank.org/"

	out, err := HttpRequest(method, url, "")
	if err != nil {
		t.Error("could not complete HttpRequest", err)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(""))
	if err != nil {
		t.Error("could not create http request", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log("http request failed", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		t.Error("could not read from buffer", err)
	}

	if !reflect.DeepEqual(buf.Bytes(), out) {
		t.Error("return value from HttpRequest does not match expected value")
	}
}

func TestCreateAuthNHeader(t *testing.T) {
	// TODO test this
}

func TestJwtHTTPRequest(t *testing.T) {
	// TODO test this
}

func TestUnrollStringSlicetoMapStringString(t *testing.T) {
	var tests = []struct {
		slices              []string
		delim               string
		expectedMap         map[int]string
		expectedStringArray []string
	}{
		{
			slices:              []string{"0/something", "1/something else", "no delim"},
			delim:               "/",
			expectedMap:         map[int]string{0: "something", 1: "something else"},
			expectedStringArray: []string{"no delim"},
		},
		{
			slices:              []string{"0/something", "1/something else"},
			delim:               "/",
			expectedMap:         map[int]string{0: "something", 1: "something else"},
			expectedStringArray: []string{},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			outMap, outString, err := UnrollStringSliceToMapIntString(tt.slices, tt.delim)
			if err != nil {
				t.Error("could not convert string slice to map[int]string", err)
			}

			if !reflect.DeepEqual(outMap, tt.expectedMap) || !reflect.DeepEqual(outString, tt.expectedStringArray) {
				t.Error("return value of UnrollStringSliceToMapIntString does not match expected value")
			}
		})
	}
}

func TestCheckLoad(t *testing.T) {
	// TODO test this
}

func TestIsTTY(t *testing.T) {
	// TODO test this
}
