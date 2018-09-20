package framework

import (
	"regexp"
	"strings"
)

func RemoveSpecialChar(char string)  string{
	reg, err := regexp.Compile("[^a-z -Z_]+")
	if err != nil {
		return ""
	}
	str := reg.ReplaceAllString(char, "")
	str = AddSlash(str)
	return str
}

func AddSlash(char string)  string{
	str := strings.Replace(char, "'", "\\'", -1)
	str = strings.Replace(char, "\"", "\\\"", -1)
	return str
}
