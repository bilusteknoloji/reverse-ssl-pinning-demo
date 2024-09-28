---
name: Bug report
about: Create a report to help us improve
title: "[BUG]"
labels: bug
assignees: vigo

---

## Describe the bug

A clear and concise description of what the bug is. Include details about the
specific HTTP request or response if relevant.

## To Reproduce

Steps to reproduce the behavior:

1. Run the server using `go run cmd/server/main.go`
1. Run the client using `go run cmd/client/main.go`
1. or use `curl` example
1. Check your certificate generation steps
1. Check provided certificates under `certs/` folder
1. If applicable, include the request and response data
1. See the error

## Expected behavior

A clear and concise description of what you expected to happen, such as
correct logging of headers or valid signature verification.

## Logs or Output

If applicable, include any error messages or output seen in the terminal/logs.

```bash
$ go run cmd/server/main.go     # for server

$ go run cmd/client/main.go     # for client

# or curl
$ curl --cacert certs/server/server-cert.pem \
     --cert certs/client/client-cert-signed.pem \
     --key certs/client/client-key.pem 
     https://localhost:8443

# Sample output or error logs here
```

## Environment (please complete the following information):

- OS: [e.g., Ubuntu, macOS]
- CPU: [e.g., M3]
- Go Version: [e.g., 1.19]
- Your SHELL and version: [e.g., bash, 5.2.32(1)-release]

