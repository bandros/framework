package framework

import (
	"bytes"
	"regexp"
	"strings"
)

func RemoveSpecialChar(char string) string {
	char = string(bytes.Trim([]byte(char), "\xef\xbb\xbf"))
	reg, err := regexp.Compile("[^ -~]+")
	if err != nil {
		return ""
	}
	str := reg.ReplaceAllString(char, "")
	str = AddSlash(str)
	return str
}

func AddSlash(char string) string {
	str := strings.Replace(char, "'", "\\'", -1)
	str = strings.Replace(char, "\"", "\\\"", -1)
	return str
}
