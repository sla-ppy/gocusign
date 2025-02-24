package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonData := `
		{
			"intValue":1234,
			"boolValue":true,
			"stringValue":"hello!",
			"dateValue":"2022-03-02T09:10:00Z",
			"objectValue":{
				"arrayValue":[1,2,3,4]
			},
			"nullStringValue":null,
			"nullIntValue":null
		}
	`
	var data map[string]interface{}
	// takes raw json and unmarshals into go variables into data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Printf("Could not unmarshal JSON: %s\n", err)
		return
	}

	fmt.Printf("JSON data: %v\n", data)

	// make sure the data we got is what we were expecting
	rawDateValue, ok := data["dateValue"] // is actually a string value
	if !ok {                              // check if value is in the map
		fmt.Printf("dateValue does not exist\n")
		return
	}
	dateValue, ok := rawDateValue.(string)
	if !ok {
		fmt.Printf("dateValue is not a string\n")
		return
	}
	fmt.Printf("Date value: %s\n", dateValue)
}
