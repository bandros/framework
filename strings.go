package framework

import (
	"bytes"
	"html"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func RemoveSpecialChar(char interface{}) interface{} {
	var reflectValue = reflect.ValueOf(char)
	var val string
	var i int
	switch reflectValue.Kind() {
	case reflect.String:
		val = strings.TrimSpace(reflectValue.String())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i = int(reflectValue.Uint())
		val = strconv.Itoa(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i = int(reflectValue.Int())
		val = strconv.Itoa(i)
	}
	val = string(bytes.Trim([]byte(val), "\xef\xbb\xbf"))
	reg, err := regexp.Compile("[^ -~]+")
	if err != nil {
		return ""
	}
	str := reg.ReplaceAllString(val, "")
	str = AddSlash(str)
	str = html.EscapeString(str)
	return str
}

func AddSlash(char string) string {
	var str = strings.Replace(char, "'", "\\'", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	return str
}

func NumberPhone(number string) string {
	reg, _ := regexp.Compile("[^0-9]+")
	number = reg.ReplaceAllString(number, "")
	reg, _ = regexp.Compile("^(08|8)")
	number = reg.ReplaceAllString(number, "628")
	return number
}
