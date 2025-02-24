# gocusign
HTTP Client which uses OpenAPI calls to communicate with an HTTP server, written in Go

Things I've learned from the project:
- HTTP network protocol, GET, POST
- JWAuth Bearer Token auth with HTTP POST is possible by addig a header to do that, aka. we send the token with each POST request
- HTTP Headers specify what information we are sending or expect to receive from server, we call them [media types](https://www.iana.org/assignments/media-types/media-types.xhtml#application) such as: application/json, application/pdf

- [Tutorials by Digitial Ocean on Golang](https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go)
- [Reqbin for HTTP theory and more](https://reqbin.com/)
- [Swagger for viewing API docs such as .yaml (F1 in VScode)](https://swagger.io/)
