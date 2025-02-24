package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type myJSON struct {
	IntValue        int       `json:"intValue"`
	BoolValue       bool      `json:"boolValue"`
	StringValue     string    `json:"stringValue"`
	DateValue       time.Time `json:"dateValue"`
	ObjectValue     *myObject `json:"objectValue"`
	NullStringValue *string   `json:"nullStringValue"`
	NullIntValue    *int      `json:"nullIntValue"`
	EmptyString     string    `json:"emptyString,omitempty"`
}

type myObject struct {
	ArrayValue []int `json:"arrayValue"`
}

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
			"nullIntValue":null,
			"extraValue":4321
		}
	`
	// any fields included in the JSON data that aren't defined on the struct are ignored by Go's JSON parser, and it will continue with the next item

	var data *myJSON
	// takes raw json and unmarshals into go variables into data
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Printf("Could not unmarshal JSON: %s\n", err)
		return
	}

	fmt.Printf("JSON struct: %v#\n", data)
	fmt.Printf("dateValue: %#v\n", data.DateValue)
	fmt.Printf("objectValue: %#v\n", data.ObjectValue)
}
