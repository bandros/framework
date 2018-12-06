package framework

import "strings"

func JoinMap(data map[string]string, sep string) string {
	result := []string{}
	for _, v := range data {
		result = append(result, v)
	}

	return strings.Join(result, sep)
}

func JoinMapKey(data map[string]interface{}, sep string) string {
	result := []string{}
	for i, _ := range data {
		result = append(result, i)
	}

	return strings.Join(result, sep)
}
