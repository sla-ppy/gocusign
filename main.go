package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type sessionInit struct {
	Company string `json:"company"` // string and int cant have value: nil, since "" and 0 are their empty values
	CaseId  string `json:"case_id"` // it can either be type or nil, we need to make it a reference
	Name    string `json:"name"`    // omitempty handles whether string is omitted or not
	Email   string `json:"email"`   // of omitempty is added and the string is empty, it wont be output on the console!
	Phone   string `json:"phone"`
}

const server_url = "https://sign-test.comnica.com"

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
		fmt.Printf("*ERR*: Could not marshal json: %s\n", err)
		return
	}
	fmt.Printf("Client: JSON data prepared:\n %s\n", jsonData)

	// send POST
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Server: %s /\n", r.Method)
			fmt.Printf("Server: query id: %s\n", r.URL.Query().Get("id"))
			// header information helps with identifying problems with the server
			fmt.Printf("Server: content-type: %s\n", r.Header.Get("content-type"))
			fmt.Printf("Server: headers:\n")
			for headerName, headerValue := range r.Header {
				fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
			}

			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("*ERR*: Server could not read request body: %s\n", err)
			}
			fmt.Printf("Server: request body: %s\n", reqBody)

			fmt.Fprintf(w, `{"message": "hello!"}`)
		})

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", serverPort),
			Handler: mux,
		}

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("*ERR*: Error running HTTP server: %s\n", err)
			}
		}
	}()

	// sleep till server is rdy after POST, get ready for GET
	time.Sleep(100 * time.Millisecond)

	// []byte since encoding/json in go returns []byte
	jsonBody := []byte(`{"client_message": "hello, server!"}`)
	// wrapper for jsonBody
	bodyReader := bytes.NewReader(jsonBody) // exists so jsonBody can be the value as the request body, since http.Request needs a value that is io.Reader

	// create HTTP request
	requestURL := fmt.Sprintf(server_url) // dont specify port, we dont care
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("*ERR*: client could not create post request: %s\n", err)
		os.Exit(1)
	}
	// Header contains information on what the content type is, this is what i'll need to set to .pdf somehow
	// application/pdf !!
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// send the actual request
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("*ERR*: client couldn't send the request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client: got response!\n")
	fmt.Printf("Client: status code: %d\n", res.StatusCode)

	// read the HTTP response body
	// ReadAll() reads from io.Reader, returns data as []byte or error
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("*ERR*: client could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Client: response body:\n %s\n", resBody)

}
