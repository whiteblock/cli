package util

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/rpc/v2/json2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

func JsonRpcCallAndPrint(method string, params interface{}) {
	reply, err := JsonRpcCall(method, params)
	switch reply.(type) {
	case string:
		_, noPretty := os.LookupEnv("NO_PRETTY")
		if noPretty {
			fmt.Println(reply.(string))
		} else {
			fmt.Printf("\033[97m%s\033[0m\n", reply.(string))
		}

		return
	}

	if err != nil {
		jsonError, ok := err.(*json2.Error)
		if ok && jsonError.Data != nil {
			res, err := json.Marshal(jsonError.Data)
			if err != nil {
				PrintErrorFatal(err)
			}
			PrintStringError(string(res))
			os.Exit(1)
		} else {
			PrintErrorFatal(err)
		}
	}
	fmt.Println(Prettypi(reply))
}
func JsonRpcCallP(method string, params interface{}, out interface{}) error {
	res, err := JsonRpcCall(method, params)
	if err != nil {
		return err
	}
	tmp, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, out)
}
func JsonRpcCall(method string, params interface{}) (interface{}, error) {
	//log.Println("URL IS "+url)
	jrpc, err := json2.EncodeClientRequest(method, params)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	body := strings.NewReader(string(jrpc))
	req, err := func() (*http.Request, error) {
		if strings.HasSuffix(conf.ServerAddr, "5000") { //5000 is http
			return http.NewRequest("POST", fmt.Sprintf("http://%s/rpc", conf.ServerAddr), body)
		} else { //5001 is https
			return http.NewRequest("POST", fmt.Sprintf("https://%s/rpc", conf.ServerAddr), body)
		}

	}()
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	auth, err := CreateAuthNHeader()
	if err != nil {
		log.Println(err)
	} else {
		req.Header.Set("Authorization", auth) //If there is an error, dont send this header for now
	}

	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	var out interface{}
	err = json2.DecodeClientResponse(resp.Body, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
