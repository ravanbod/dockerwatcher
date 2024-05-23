package jsontotree

import (
	"encoding/json"
)

func ConvertJsonToTree(jsonData string) string {
	var result string
	var theMap map[string]interface{}

	json.Unmarshal([]byte(jsonData), &theMap)

	for k, v := range theMap {
		stringVar, ok := v.(string)
		if ok {
			result += k + ": " + stringVar + "\n"
		} else {
			innerMap, ok := v.(map[string]interface{})
			if ok {
				innerMapJson, err := json.Marshal(innerMap)
				if err == nil {
					result += ConvertJsonToTree(string(innerMapJson))
				}
			}
		}
	}
	return result
}
