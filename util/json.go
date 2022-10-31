package util

import "encoding/json"

func ToJson(v interface{}) string {
	result, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(result)
}
