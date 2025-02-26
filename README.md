# gocusign
HTTP Client which uses OpenAPI calls to communicate with an HTTP server, written in Go

Things I've learned from the project:
- HTTP network protocol, GET, POST
- JWAuth Bearer Token auth with HTTP POST is possible by addig a header to do that, aka. we send the token with each POST request
- HTTP Headers specify what information we are sending or expect to receive from server, we call them [media types](https://www.iana.org/assignments/media-types/media-types.xhtml#application) such as: application/json, application/pdf
- Marshal from Go data to multipart/form-data
- Marshal from Go data to JSON, unmarshaling JSON to Go

## Testing HTTP requests locally using netcat and curl
- Set baseUrl to localhost addr
- Even when the program crashes the requests till that point can be examined
- Important: Only HTTP works! HTTPS doesn't work

Terminal #0
```
nc -l -p 1234 -v | less     -> nc/netcat [TCP/IP server testing tool], -l [listens on the port], -p [sets port], -v [return more info]
```
Terminal #1
```
curl http://127.0.0.1:1234/     -> make HTTP GET request to test
```

- [Tutorials by Digitial Ocean on Golang](https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go)
- [Reqbin for HTTP theory and more](https://reqbin.com/)
- [Swagger for viewing API docs such as .yaml (F1 in VScode)](https://swagger.io/)
