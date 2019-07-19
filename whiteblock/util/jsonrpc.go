package util

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/rpc/v2/json2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func JsonRpcCallAndPrint(method string, params interface{}) {
	reply, err := JsonRpcCall(method, params)
	if err != nil {
		jsonError, ok := err.(*json2.Error)
		if ok && jsonError.Data != nil {
			res, err := json.Marshal(jsonError.Data)
			if err != nil {
				PrintErrorFatal(err)
			}
			PrintErrorFatal(string(res))
		} else {
			PrintErrorFatal(err)
		}
	}
	Print(reply)
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
	var res interface{}
	var err error
	for i := 0; i < conf.RPCRetries; i++ {
		res, err = jsonRpcCall(method, params)
		if err == nil {
			break
		}
	}
	return res, err
}

func jsonRpcCall(method string, params interface{}) (interface{}, error) {
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
