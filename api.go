package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type Api struct {
	Url         string
	Data        map[string]interface{}
	dataJsonRaw string
	Header      map[string]string
	ContentType string
	Username    string
	Password    string
	BasicAuth   bool
	body        string
}

func (api *Api) Do(method string) error {
	method = strings.ToUpper(method)
	client := &http.Client{}
	var err error
	var req *http.Request
	if api.dataJsonRaw != "" {
		var reader = strings.NewReader(api.dataJsonRaw)
		req, err = http.NewRequest(method, api.Url, reader)
		if err != nil {
			return err
		}
	} else if method == "POST" {
		param := url.Values{}
		for i, v := range api.Data {
			var reflectValue = reflect.ValueOf(v)
			switch reflectValue.Kind() {
			case reflect.String:
				param.Set(i, v.(string))
			case reflect.Slice:
				if reflect.TypeOf(v).String() != "[]string" {
					return errors.New("slice only support []string type")
				}
				for _, v2 := range v.([]string) {
					var index = i + "[]"
					param.Add(index, v2)
				}

			case reflect.Map:
				if reflect.TypeOf(v).String() != "map[string]string" {
					return errors.New("map only support map[string]string type")
				}
				for i2, v2 := range v.(map[string]string) {
					var index = i + "[" + i2 + "]"
					param.Add(index, v2)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var str = strconv.Itoa(int(reflectValue.Uint()))
				param.Set(i, str)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var str = strconv.Itoa(int(reflectValue.Int()))
				param.Set(i, str)
			default:
				return errors.New(reflectValue.String() + " Not support")
			}
		}
		req, err = http.NewRequest(method, api.Url, bytes.NewBufferString(param.Encode()))
		if api.BasicAuth {
			req.SetBasicAuth(api.Username, api.Password)
		}
		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequest("GET", api.Url, nil)
		if api.BasicAuth {
			req.SetBasicAuth(api.Username, api.Password)
		}
		if err != nil {
			return err
		}

		param := req.URL.Query()
		for i, v := range api.Data {
			if reflect.TypeOf(v).String() != "string" {
				return errors.New("only support string type")
			}
			param.Set(i, v.(string))
		}
		req.URL.RawQuery = param.Encode()
	}

	///api.Header = map[string]string{}
	if api.Header == nil {
		api.Header = map[string]string{}
	}
	if api.ContentType == "" {
		api.Header["Content-Type"] = "application/x-www-form-urlencoded"
	} else {
		api.Header["Content-Type"] = api.ContentType
	}
	for i, v := range api.Header {
		req.Header.Set(i, v)
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

func (api *Api) Get(data interface{}) error {
	err := json.Unmarshal([]byte(api.body), &data)
	if err != nil {
		return err
	}
	return nil
}

func (api *Api) GetXml(data interface{}) error {
	err := xml.Unmarshal([]byte(api.body), &data)
	if err != nil {
		return err
	}
	return nil
}

func (api *Api) GetRaw() string {
	return api.body
}

func (api *Api) JsonData(raw interface{}) error {
	var jsonByte, err = json.Marshal(raw)
	if err != nil {
		return err
	}
	api.dataJsonRaw = string(jsonByte)
	return nil
}
