package main

import (
	"encoding/json"
	"fmt"
	"time"
)

/* names must start with a capital letter
// struct tags define how the json is generated

type myInt struct {
	IntValue int
}
=> output: {"IntValue":1234}
*/

/*
type myInt struct {
	IntValue int `json:"intValue"`
}
=> output: {"intValue":1234}
*/

type myJSON struct {
	IntValue        int       `json:"intValue"`
	BoolValue       bool      `json:"boolValue"`
	StringValue     string    `json:"stringValue"`
	DateValue       time.Time `json:"datetValue"`
	ObjectValue     *myObject `json:"objectValue"`
	NullStringValue *string   `json:"nullStringValue,omitempty"` // string and int cant have value: nil, since "" and 0 are their empty values
	NullIntValue    *int      `json:"nullIntValue"`              // it can either be type or nil, we need to make it a reference
	EmptyString     string    `json:"emptyString,omitempty"`     // omitempty handles whether string is omitted or not
} // of omitempty is added and the string is empty, it wont be output on the console!

type myObject struct {
	ArrayValue []int `json:"arrayValue"`
}

func main() {
	otherInt := 4321
	data := &myJSON{
		IntValue:    1234,
		BoolValue:   true,
		StringValue: "hello!",
		DateValue:   time.Date(2022, 3, 2, 9, 10, 0, 0, time.UTC),
		ObjectValue: &myObject{
			ArrayValue: []int{1, 2, 3, 4},
		},
		NullStringValue: nil,
		NullIntValue:    &otherInt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	fmt.Printf("json data: %s\n", jsonData)
}
