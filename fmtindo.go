package framework

import (
	"fmt"
	"strconv"
	"strings"
)

func DateTime(text string, format string) string {
	result := format
	split := strings.Split(text, " ")
	date := strings.Split(split[0], "-")
	//time := strings.Split(split[1],":")

	//Year
	result = strings.Replace(result, "&Y", date[0], -1) // 2018
	result = strings.Replace(result, "&y", date[0], -1) // 18

	//Month
	monthNumber, err := strconv.ParseUint(date[1], 10, 8)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	m := uint8(monthNumber)
	result = strings.Replace(result, "&M", monthID(m, false), -1)
	result = strings.Replace(result, "&m", monthID(m, true), -1)

	//Day
	result = strings.Replace(result, "&D", date[2], -1)
	result = strings.Replace(result, "&d", date[2], -1)

	return result

}

func monthID(month uint8, shortMonth bool) string {
	shortMonthID := []string{
		" ",
		"Jan",
		"Feb",
		"Mar",
		"Aprl",
		"Mei",
		"Jun",
		"Jul",
		"Agst",
		"Sept",
		"Okt",
		"Nov",
		"Des",
	}

	longMonthID := []string{
		" ",
		"Januari",
		"Februari",
		"Maret",
		"April",
		"Mei",
		"Juni",
		"Juli",
		"Agustus",
		"September",
		"Oktber",
		"Novmeber",
		"Desember",
	}

	if shortMonth {
		return shortMonthID[month]
	}

	return longMonthID[month]
}
