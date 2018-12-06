package framework

import (
	"regexp"
	"strconv"
)

func SplitN(String string, n int) []string {
	number := strconv.Itoa(n)
	re := regexp.MustCompile(".{0," + number + "}")
	data := re.FindAllString(String, -1)
	return data
}
