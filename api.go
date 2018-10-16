package framework

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Api struct {
	Url string
	Data map[string]string
	Header map[string]string
	ContentType string
	body string
}

func(api *Api) Do(method string)  error{
	method = strings.ToUpper(method)
	client := &http.Client{}
	var err error
	var req *http.Request
	if method == "POST" {
		param := url.Values{}
		for i,v := range api.Data {
			param.Set(i,v)
		}
		req, err = http.NewRequest(method, api.Url, bytes.NewBufferString(param.Encode()))
		if err != nil {
			return err
		}
	}else{
		req, err = http.NewRequest("GET", api.Url, nil)
		if err != nil {
			return err
		}
		param := req.URL.Query()
		for i,v := range api.Data {
			param.Set(i,v)
		}
		req.URL.RawQuery = param.Encode()
	}

	///api.Header = map[string]string{}
	if api.Header == nil {
		api.Header = map[string]string{}
	}
	if api.ContentType == "" {
		api.Header["Content-Type"] = "application/x-www-form-urlencoded"
	}else{
		api.Header["Content-Type"] = api.ContentType
	}
	for i,v := range api.Header {
		req.Header.Set(i,v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	api.body = string(body)
	return nil
}

func(api *Api) Get(data interface{}) (error){
	err := json.Unmarshal([]byte(api.body), &data)
	if err != nil {
		return err
	}
	return nil
}

func(api *Api) GetRaw() (string){
	return api.body
}
