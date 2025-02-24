package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type sessionInit struct {
	Company string `json:"company"` // string and int cant have value: nil, since "" and 0 are their empty values
	CaseId  string `json:"case_id"` // it can either be type or nil, we need to make it a reference
	Name    string `json:"name"`    // omitempty handles whether string is omitted or not
	Email   string `json:"email"`   // of omitempty is added and the string is empty, it wont be output on the console!
	Phone   string `json:"phone"`
}

const baseUrl = "https://sign-test.comnica.com/api/session"
const authToken = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiY3JlYXRlX3Nlc3Npb24iXSwiaXNzIjoic2lnbi1iYWNrZW5kLXRlc3QiLCJleHAiOjE3NDM4NzQzNDYsImlhdCI6MTczOTU1NDM0Niwic3ViIjoiY29tcGFueTpiYWNrZW5kLWRldmVsb3Blci5jb21uaWNhLmlkIn0.x9qB3JtDtl-cGp9ijyjaL-lVuRSQsOk-KVibU8p3eyk"

func main() {
	// 0. initialize data to serialize
	data := &sessionInit{
		Company: "backend-developer.comnica.id",
		CaseId:  "sla-ppy",
		Name:    "Miklos Vida",
		Email:   "jackyissocial@gmail.com",
		Phone:   "36203925891",
	}
	
	// 1. serialize from Go data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("*ERR*: Could not marshal JSON: %s\n", err)
		return
	}
	fmt.Printf("Client: JSON data prepared:\n%s\n", jsonData)

	// prepare post for /session/init
	req, err := http.NewRequest("POST", baseUrl+"/init", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")   // tell the server we are sending json data
	req.Header.Add("Accept", "application/json")         // accept tells the server we are expecting json back
	req.Header.Add("Authorization", "Bearer "+authToken) // auth token is a must for all post requests for us

	// actually send post
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("*ERR*: Establishing /session/init failed! Status: ", resp.Status)
		panic(err)
	}
	defer resp.Body.Close()

	// GET response
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	fmt.Println("Response Body:", resp.Body)
	fmt.Println("Response Body:", resp)

	const newAuthToken = ""
	const sessionId = ""
}
