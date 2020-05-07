package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
)

type Apis struct {
	Url         string
	Data        map[string]interface{}
	Header      map[string]string
	ContentType string
	Username    string
	Password    string
	BasicAuth   bool
	body        string
	Status      string
}

func (api *Apis) Do(method string) error {
	url := api.Url
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	for i, v := range api.Data {
		var reflectValue = reflect.ValueOf(v)
		switch dt := v.(type) {
		case string:
			writer.WriteField(i, dt)
		case []string:
			for _, v2 := range dt {
				var index = i + "[]"
				writer.WriteField(index, v2)
			}
		case map[string]string:
			for i2, v2 := range dt {
				var index = i + "[" + i2 + "]"
				writer.WriteField(index, v2)
			}
		case uint, uint8, uint16, uint32, uint64:
			var str = strconv.Itoa(int(reflectValue.Uint()))
			writer.WriteField(i, str)
		case int, int8, int16, int32, int64:
			var str = strconv.Itoa(int(reflectValue.Int()))
			writer.WriteField(i, str)
		case *multipart.FileHeader:
			file, _ := dt.Open()
			write, _ := writer.CreateFormFile(i, dt.Filename)
			fmt.Println(dt.Filename)
			_, _ = io.Copy(write, file)
			_ = file.Close()
		case []*multipart.FileHeader:
			for _, v2 := range dt {
				file, _ := v2.Open()
				f, _ := writer.CreateFormFile(i+"[]", v2.Filename)
				_, _ = io.Copy(f, file)
				_ = file.Close()
			}
		default:
			return errors.New(reflectValue.String() + " Not support")
		}
	}

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}

	if api.Header == nil {
		api.Header = map[string]string{}
	}
	if api.ContentType == "" {
		api.Header["Content-Type"] = writer.FormDataContentType()
	} else {
		api.Header["Content-Type"] = api.ContentType
	}

	for i, v := range api.Header {
		req.Header.Set(i, v)
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	api.body = string(body)
	api.Status = res.Status

	return nil
}

func (api *Apis) Get(data interface{}) error {
	err := json.Unmarshal([]byte(api.body), &data)
	if err != nil {
		return err
	}
	return nil
}

func (api *Apis) GetXml(data interface{}) error {
	err := xml.Unmarshal([]byte(api.body), &data)
	if err != nil {
		return err
	}
	return nil
}

func (api *Apis) GetRaw() string {
	return api.body
}
