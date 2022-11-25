package jasonconverter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func convertJSON(data []byte) []byte {
	parsed := map[string]interface{}{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		fmt.Println(err)
	}

	result, err := json.Marshal(camel2snake(parsed))
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func camel2snake(data map[string]interface{}) (result map[string]interface{}) {
	result = map[string]interface{}{}

	converter := regexp.MustCompile("([a-z0-0])([A-Z])")
	for k, v := range data {
		snake := strings.ToLower(converter.ReplaceAllString(k, "${1}_${2}"))

		switch v.(type) {
		case []interface{}:
			temp := []interface{}{}

			for _, obj := range v.([]interface{}) {
				switch obj.(type) {
				case map[string]interface{}:
					temp = append(temp, camel2snake(obj.(map[string]interface{})))
				default:
					temp = append(temp, obj)
				}
			}
			result[snake] = temp

		case map[string]interface{}:
			temp := camel2snake(v.(map[string]interface{}))
			result[snake] = temp

		default:
			result[snake] = v
		}
	}

	return result
}
