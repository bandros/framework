package framework

import "reflect"

func InArray(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func InMap(mapData []map[string]interface{}, key string, item interface{}) bool {
	for _, v := range mapData {
		if v[key] == item {
			return true
		}
	}

	return false
}
