package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// string translates directly into JSON object key
	// interface{} - allows any data type
	data := map[string]interface{}{
		"intValue":    1234,
		"boolValue":   true,
		"stringValue": "hello!",
		"objectValue": map[string]interface{}{
			"arrayValue": []int{1, 2, 3, 4},
		},
		"dateValue":       time.Date(2022, 3, 2, 9, 10, 0, 0, time.UTC),
		"nullStringValue": nil,
		"nullIntValue":    nil,
	}

	// Marshalling == Serialization
	// Serialization: Transformation of program data in memory to an easily transferable format
	// Ex.: Go data into JSON
	jsonData, err := json.Marshal(data) // Marshal() automatically decides the type for an interface{} object
	if err != nil {
		fmt.Printf("Could not marshal json: %s\n", err)
		return
	}

	fmt.Printf("Json data: %s\n", jsonData)
}
