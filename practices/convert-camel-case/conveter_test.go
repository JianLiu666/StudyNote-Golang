package convertcamelcase

import (
	"fmt"
	"testing"
)

func TestConverter(t *testing.T) {
	str := `{
		"stringField": "ShouldNotBeChanged",
		"numberField": 0,
		"stringArray": ["ShouldNotBeChanged", "ShouldNotBeChanged"],
		"numberArray": [0, 1],
		"nestedObj1": [
			{
				"fieldName11": 0,
				"fieldName12": "ShouldNotBeChanged"
			}
		],
		"nestedObj2": {
			"fieldName21": 0,
			"fieldName22": "ShouldNotBeChanged",
			"nestedObj23": [
				{
					"fieldName231": 0,
					"fieldName232": "ShouldNotBeChanged"
				}
			]
		},
		"nestedObj3": {
			"nestedObj31": {
				"nestedObj32": {
					"fieldName33": "ShouldNotBeChanged"
				}
			}
		}
	}`

	result := convertJSON([]byte(str))
	fmt.Println(string(result))
}
