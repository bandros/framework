package framework

import (
	"fmt"
	"strconv"
	"strings"
)

func NumberFormat(d int64) string {
	var sep = "."
	var sign = ""
	if d < 0 {
		sign = "-"
		d *= -1
	}

	var path = []string{"", "", "", "", "", ""}
	var i = len(path) - 1

	for d > 999 {
		//get mod 1000 and to string also change to 3 digit
		path[i] = fmt.Sprintf("%03d", d%1000)
		d /= 1000
		i--
	}
	path[i] = strconv.FormatInt(d, 10)
	return sign + strings.Join(path[i:], sep)
}
