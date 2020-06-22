package utils

import (
	"reflect"
)

func GetMapValue(key string, msgMap map[string]interface{}) (value string) {
	if msgMap[key] == nil {
		value = ""
	} else {
		value = msgMap[key].(string)
	}
	return value
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}