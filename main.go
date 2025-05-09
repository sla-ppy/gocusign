package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

const baseUrl = "https://sign-test.comnica.com/api/session"

//const baseUrl = "http://127.0.0.1:1234/"

func sessionInit() (string, string) {
	// define JSON data structure
	type initRequest struct {
		Company string `json:"company"` // string and int cant have value: nil, since "" and 0 are their empty values
		CaseId  string `json:"case_id"` // it can either be type or nil, we need to make it a reference
		Name    string `json:"name"`    // omitempty handles whether string is omitted or not
		Email   string `json:"email"`   // of omitempty is added and the string is empty, it won't be output on the console!
		Phone   string `json:"phone"`
	}

	// initialize data to serialize
	requestData := &initRequest{
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

	// prepare post for /session/init
	const authToken = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiY3JlYXRlX3Nlc3Npb24iXSwiaXNzIjoic2lnbi1iYWNrZW5kLXRlc3QiLCJleHAiOjE3NDM4NzQzNDYsImlhdCI6MTczOTU1NDM0Niwic3ViIjoiY29tcGFueTpiYWNrZW5kLWRldmVsb3Blci5jb21uaWNhLmlkIn0.x9qB3JtDtl-cGp9ijyjaL-lVuRSQsOk-KVibU8p3eyk"
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("*ERR*: Closing the body failed: %s\n", err)
			panic(err)
		}
	}(resp.Body)

	// check if OK, proceed with unmarshalling
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Successful /session/init call! Status Code is OK: %d\n", resp.StatusCode)
	} else {
		fmt.Printf("*ERR*: Status Code is not OK: %d\n", resp.StatusCode)
		panic(err)
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
		panic(err)
	}

	// expected server response in json format
	var respData struct {
		SessionId   string `json:"session_id"`
		BearerToken string `json:"bearer_token"`
	}

	// unmarshal from go to json
	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Printf("*ERR*: Unmarshalling was unsuccessful: %s\n", err)
		panic(err)
	}

	fmt.Printf("{\"session_id\":\"%s\",\"bearer_token\":\"%s\"}\n", respData.SessionId, respData.BearerToken)

	return respData.SessionId, respData.BearerToken
}

/*
since multipart.writer doesn't support specifying content type, and curl returned the following result for sending .pdf data:
-F 'data=@input.pdf;type=application/pdf' \
we need this function to convert .pdf properly
*/
func setPdfContentType(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "application/pdf")
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "data", filename))
	return w.CreatePart(h)
}

func sessionAddDocument(sessionId string, bearerToken string) {
	// multipart/form-data initialization
	formData := &bytes.Buffer{} // we don't need to serialize here, since we already have a stream of bytes
	writer := multipart.NewWriter(formData)
	writer.WriteField("session_id", sessionId)
	writer.WriteField("description", "Example description")
	writer.WriteField("filename", "output.pdf")
	writer.WriteField("document_type", "user_document")
	fileWriter, err := setPdfContentType(writer, "vevokeszulek.pdf")
	if err != nil {
		fmt.Printf("*ERR*: Could not set .pdf content type: %s\n", err)
		panic(err)
	}
	// check if file can be opened
	fileReader, err := os.Open("resources/vevokeszulek.pdf")
	if err != nil {
		fmt.Printf("*ERR*: Could not open .pdf file: %s\n", err)
		panic(err)
	} else {
		fmt.Printf("Client: File reading succesful!\n")
	}
	// copy
	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		fmt.Printf("*ERR*: Could not copy Writer: %s\n", err)
		panic(err)
	}
	writer.Close()

	// prepare post for /session/add_document
	req, err := http.NewRequest("POST", baseUrl+"/add_document", formData)
	req.Header.Set("Content-Type", writer.FormDataContentType()) // tell the server we are sending json data
	req.Header.Set("Accept", "application/json")                 // accept tells the server we are expecting json back
	req.Header.Set("Authorization", "Bearer "+bearerToken)       // auth token is a must for all post requests for us

	// make POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("*ERR*: Establishing /session/add_document failed! Status: %s\n", err)
		panic(err)
	}
	defer func(Body io.ReadCloser) { // defer = destructor in terms of working
		err := Body.Close()
		if err != nil {
			fmt.Printf("*ERR*: Closing the body failed: %s\n", err)
			panic(err)
		}
	}(resp.Body)

	// check if OK, proceed with unmarshalling
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Succesful /session/add_document call! Status Code is OK: %d\n", resp.StatusCode)
	} else {
		fmt.Printf("*ERR*: Status Code is not OK: %d\n", resp.StatusCode)
		// don't panic, we might get additional info back to debug
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
		panic(err)
	}

	// expected server response in json format
	var respData struct {
		DocumentId int `json:"document_id"`
	}

	// unmarshal from go to json
	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Printf("*ERR*: Unmarshalling was unsuccesful: %s\n", err)
		panic(err)
	}

	fmt.Printf("{\"document_id\":\"%d\"}\n", respData.DocumentId)
}

func sessionCheckState(sessionId string, bearerToken string) string {
	type initRequest struct {
		SessionId string `json:"session_id"`
	}

	// initialize data to serialize
	requestData := &initRequest{
		SessionId: sessionId,
	}

	// serialize from Go data to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("*ERR*: Could not marshal JSON: %s\n", err)
		panic(err)
	}
	fmt.Printf("Client: JSON data prepared:\n%s\n", jsonData)

	// prepare post for /session/check_state
	req, err := http.NewRequest("POST", baseUrl+"/check_state", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")     // tell the server we are sending json data
	req.Header.Add("Accept", "application/json")           // accept tells the server we are expecting json back
	req.Header.Add("Authorization", "Bearer "+bearerToken) // auth token is a must for all post requests for us

	// make POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("*ERR*: Establishing /session/init failed! Status: ", resp.Status)
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("*ERR*: Closing the body failed: %s\n", err)
			panic(err)
		}
	}(resp.Body)

	// check if OK, proceed with unmarshalling
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Succesful /session/check_state call! Status Code is OK: %d\n", resp.StatusCode)
	} else {
		fmt.Printf("*ERR*: Status Code is not OK: %d\n", resp.StatusCode)
		panic(err)
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
		panic(err)
	}

	// expected server response in json format
	var respData struct {
		State string `json:"state"`
	}

	// unmarshal from go to json
	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Printf("*ERR*: Unmarshalling was unsuccessful: %s\n", err)
		panic(err)
	}

	return respData.State
}

func sessionReady(sessionId string, bearerToken string) string {
	type initRequest struct {
		SessionId string `json:"session_id"`
	}

	// initialize data to serialize
	requestData := &initRequest{
		SessionId: sessionId,
	}

	// serialize from Go data to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("*ERR*: Could not marshal JSON: %s\n", err)
		panic(err)
	}
	fmt.Printf("Client: JSON data prepared:\n%s\n", jsonData)

	// prepare post for /session/ready
	req, err := http.NewRequest("POST", baseUrl+"/ready", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")     // tell the server we are sending json data
	req.Header.Add("Accept", "application/json")           // accept tells the server we are expecting json back
	req.Header.Add("Authorization", "Bearer "+bearerToken) // auth token is a must for all post requests for us

	// make POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("*ERR*: Establishing /session/ready failed! Status: ", resp.Status)
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("*ERR*: Closing the body failed: %s\n", err)
			panic(err)
		}
	}(resp.Body)

	// check if OK, proceed with unmarshalling
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Succesful /session/ready call! Status Code is OK: %d\n", resp.StatusCode)
	} else {
		fmt.Printf("*ERR*: Status Code is not OK: %d\n", resp.StatusCode)
		panic(err)
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
		panic(err)
	}

	// expected server response in json format
	var respData struct {
		SessionId string `json:"session_id"`
	}

	// unmarshal from go to json
	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Printf("*ERR*: Unmarshalling was unsuccessful: %s\n", err)
		panic(err)
	}

	return respData.SessionId
}

func main() {
	// 1. Init session
	sessionId, bearerToken := sessionInit()
	fmt.Printf("\n")

	state := sessionCheckState(sessionId, bearerToken)
	if state != "started" {
		fmt.Printf("*ERR*: Session state is incorrect, expected state: started, result: %s\n", state)
		panic(1)
	} else {
		fmt.Printf("State: %s\n", state)
		fmt.Printf("\n")
	}

	// 2. Add documents
	sessionAddDocument(sessionId, bearerToken)
	fmt.Printf("\n")

	state = sessionCheckState(sessionId, bearerToken)
	if state != "started" {
		fmt.Printf("*ERR*: Session state is incorrect, expected state: documents_added, result: %s\n", state)
		panic(1)
	} else {
		fmt.Printf("State: %s\n", state)
		fmt.Printf("\n")
	}

	// 3. Ready session
	sessionReady(sessionId, bearerToken)
	fmt.Printf("\n")

	state = sessionCheckState(sessionId, bearerToken)
	if state != "documents_added" {
		fmt.Printf("*ERR*: Session state is incorrect, expected state: documents_added, result: %s\n", state)
		panic(1)
	} else {
		fmt.Printf("State: %s\n", state)
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	fmt.Printf("---Program complete---\n")
	fmt.Printf("Session ID: %s\n", sessionId)
	fmt.Printf("Link for the user: %s\n", "https://sign-test.comnica.com/"+sessionId)
	fmt.Printf("---Program complete---\n")
}
