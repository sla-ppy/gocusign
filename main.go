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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("*ERR*: Closing the body failed: %s\n", err)
			panic(err)
		}
	}(resp.Body)

	// expected server response in json format
	var initResult struct {
		SessionId   string `json:"session_id"`
		BearerToken string `json:"bearer_token"`
	}

	// check if OK, proceed with unmarshalling
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Successful /session/init call! Status Code is OK: %d\n", resp.StatusCode)

		// read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
			panic(err)
		}

		// unmarshal from go to json
		err = json.Unmarshal(body, &initResult)
		if err != nil {
			fmt.Printf("*ERR*: Unmarshalling was unsuccessful: %s\n", err)
			panic(err)
		}

		fmt.Printf("{\"session_id\":\"%s\",\"bearer_token\":\"%s\"}\n", initResult.SessionId, initResult.BearerToken)
	} else {
		fmt.Printf("*ERR*: Status Code is not OK: %d\n", resp.StatusCode)
		panic(err)
	}
	return initResult.SessionId, initResult.BearerToken
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

func sessionAddDocument(sessionId string, bearerToken string) int {
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
	// copy?
	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		fmt.Printf("*ERR*: Could not copy Writer: %s\n", err)
		panic(err)
	}
	writer.Close()

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
		// dont panic, we might get additional info back to debug
	}

	// expected server response in json format
	var addDocumentResult struct {
		DocumentId int `json:"document_id"`
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("*ERR*: Cannot read response body: %s\n", err)
		panic(err)
	}

	// unmarshal from go to json
	err = json.Unmarshal(body, &addDocumentResult)
	if err != nil {
		fmt.Printf("*ERR*: Unmarshalling was unsuccesful: %s\n", err)
		panic(err)
	}

	return addDocumentResult.DocumentId
}

func sessionCheckState() (string, string) {
	return "", ""
}

func sessionReady() (string, string) {
	return "", ""
}

func main() {
	sessionId, bearerToken := sessionInit()
	documentId := sessionAddDocument(sessionId, bearerToken)
	//documentId := sessionAddDocument("hello1", "hello2")

	fmt.Printf("Document ID: %d\n", documentId)
}
