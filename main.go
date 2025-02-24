package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const serverPort = 3333

func main() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("server: %s /\n", r.Method)
			fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
			// header information helps with identifying problems with the server
			fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type"))
			fmt.Printf("server: headers:\n")
			for headerName, headerValue := range r.Header {
				fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
			}

			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("server: could not read request body: %s\n", err)
			}
			fmt.Printf("server: request body: %s\n", reqBody)

			fmt.Fprintf(w, `{"message": "hello!"}`)
		})

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", serverPort),
			Handler: mux,
		}

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("error running http server: %s\n", err)
			}
		}
	}()

	// sleep till server is rdy
	time.Sleep(100 * time.Millisecond)

	// []byte since encoding/json in go returns []byte
	jsonBody := []byte(`{"client_message": "hello, server!"}`)
	// wrapper for jsonBody
	bodyReader := bytes.NewReader(jsonBody) // exists so jsonBody can be the value as the request body, since http.Request needs a value that is io.Reader

	// create HTTP request
	requestURL := fmt.Sprintf("http://localhost:%d?id=1234", serverPort)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	// send the actual request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	// read the HTTP response body
	// ReadAll() reads from io.Reader, returns data as []byte or error
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
}
