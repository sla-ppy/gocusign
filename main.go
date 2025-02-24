package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://sign-test.comnica.com/api/session"

func sessionInit() (string, string) {
	// define JSON data structure
	type sessionInitRequest struct {
		Company string `json:"company"` // string and int cant have value: nil, since "" and 0 are their empty values
		CaseId  string `json:"case_id"` // it can either be type or nil, we need to make it a reference
		Name    string `json:"name"`    // omitempty handles whether string is omitted or not
		Email   string `json:"email"`   // of omitempty is added and the string is empty, it won't be output on the console!
		Phone   string `json:"phone"`
	}

	// initialize data to serialize
	requestData := &sessionInitRequest{
		Company: "backend-developer.comnica.id",
		CaseId:  "sla-ppy",
		Name:    "Miklos Vida",
		Email:   "jackyissocial@gmail.com",
		Phone:   "36203925891",
	}

	// serialize from Go data to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("*ERR*: Could not marshal JSON: %s\n", err)
		panic(err)
	}
	fmt.Printf("Client: JSON data prepared:\n%s\n", jsonData)

	const authToken = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiY3JlYXRlX3Nlc3Npb24iXSwiaXNzIjoic2lnbi1iYWNrZW5kLXRlc3QiLCJleHAiOjE3NDM4NzQzNDYsImlhdCI6MTczOTU1NDM0Niwic3ViIjoiY29tcGFueTpiYWNrZW5kLWRldmVsb3Blci5jb21uaWNhLmlkIn0.x9qB3JtDtl-cGp9ijyjaL-lVuRSQsOk-KVibU8p3eyk"
	// prepare post for /session/init
	req, err := http.NewRequest("POST", baseUrl+"/init", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")   // tell the server we are sending json data
	req.Header.Add("Accept", "application/json")         // accept tells the server we are expecting json back
	req.Header.Add("Authorization", "Bearer "+authToken) // auth token is a must for all post requests for us

	// make POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("*ERR*: Establishing /session/init failed! Status: ", resp.Status)
		panic(err)
	}
	defer resp.Body.Close()

	// expected server response in json format
	var initResult struct {
		SessionId   string `json:"session_id"`
		BearerToken string `json:"bearer_token"`
	}

	// check if OK, proceed with unmarshalling
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Succesful /session/init call! Status Code is OK: %d\n", resp.StatusCode)

		// read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
			panic(err)
		}

		// unmarshal from go to json
		err = json.Unmarshal(body, &initResult)
		if err != nil {
			fmt.Printf("*ERR*: Unmarshalling was unsuccesful: %s\n", err)
			panic(err)
		}

		fmt.Printf("{\"session_id\":\"%s\",\"bearer_token\":\"%s\"}\n", initResult.SessionId, initResult.BearerToken)
	} else {
		fmt.Printf("*ERR*: Status Code is not OK: %d\n", resp.StatusCode)
		panic(err)
	}
	return initResult.SessionId, initResult.BearerToken
}

func main() {
	sessionId, bearerToken := sessionInit()

	fmt.Println("Session ID:", sessionId)
	fmt.Println("Bearer Token:", bearerToken)

	// http.Response fields i can use
	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	//fmt.Println("Response Body:", resp.Body)
	//fmt.Println("Response Body:", resp)
}
