![Version](https://img.shields.io/badge/version-0.0.1-orange.svg)
![Go](https://img.shields.io/github/go-mod/go-version/bilusteknoloji/reverse-ssl-pinning-demo)
[![Golang CI Lint](https://github.com/bilusteknoloji/reverse-ssl-pinning-demo/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/bilusteknoloji/reverse-ssl-pinning-demo/actions/workflows/golangci-lint.yml)
![Powered by Rake](https://img.shields.io/badge/powered_by-rake-blue?logo=ruby)

# Reverse SSL Pinning Demo

Reverse SSL pinning is a process where the server **validates** that the
client is using the **expected certificate** or **key**. This technique can be
used to **ensure** that **only certain trusted clients** or **applications**
are **allowed** to communicate with the server.

This is a super basic approach to ssl reverse pinning approach.

---

## Requirements

Example certificates can be found under `certs/` folder:

```bash
$ tree certs/

certs/
├── client
│   ├── client-ca-key.pem        # CA Private Key
│   ├── client-ca.pem            # CA Certificate
│   ├── client-ca.srl            # CA Serial Number File
│   ├── client-cert-signed.pem   # Private key to sign the client certificate
│   ├── client-cert.csr          # Client Certificate Signing Request
│   └── client-key.pem           # Client Key
└── server
    ├── server-cert.csr          # Server Certificate Signing Request
    ├── server-cert.pem          # Server Certificate
    └── server-key.pem           # Server Key
```

You can create your own keys and certificates!

```bash
# CA certificate
openssl req -new -newkey rsa:2048 -nodes -x509 -days 365 -keyout certs/client/client-ca-key.pem -out certs/client/client-ca.pem -config ca_openssl.cnf

# Server Key and CSR
openssl req -new -newkey rsa:2048 -nodes -keyout certs/server/server-key.pem -out certs/server/server-cert.csr -config openssl.cnf

# Sign the Server CSR with CA
openssl x509 -req -in certs/server/server-cert.csr -CA certs/client/client-ca.pem -CAkey certs/client/client-ca-key.pem -CAcreateserial -out certs/server/server-cert.pem -days 365 -extensions req_ext -extfile openssl.cnf
# Certificate request self-signature ok
# subject=C=TR, ST=Istanbul, L=Istanbul, O=Bilus, OU=Development, CN=localhost

# Client Key and CSR
openssl req -new -newkey rsa:2048 -nodes -keyout certs/client/client-key.pem -out certs/client/client-cert.csr -config openssl.cnf

# Sign the Client CSR with CA
openssl x509 -req -in certs/client/client-cert.csr -CA certs/client/client-ca.pem -CAkey certs/client/client-ca-key.pem -CAcreateserial -out certs/client/client-cert-signed.pem -days 365 -extensions req_ext -extfile openssl.cnf
# Certificate request self-signature ok
# subject=C=TR, ST=Istanbul, L=Istanbul, O=Bilus, OU=Development, CN=localhost
```

Now check your certificate:

```bash
openssl x509 -in certs/server/server-cert.pem -text -noout
```

You should see (*partial*):

    Certificate:
        Data:
            :
            Signature Algorithm: sha256WithRSAEncryption
            Issuer: C=TR, ST=Istanbul, L=Istanbul, O=Bilus, OU=Development, CN=Bilus Root CA
            Validity
                Not Before: Sep 28 16:54:33 2024 GMT
                Not After : Sep 28 16:54:33 2025 GMT
            Subject: C=TR, ST=Istanbul, L=Istanbul, O=Bilus, OU=Development, CN=localhost
            :
            :
            X509v3 extensions:
                X509v3 Subject Alternative Name: 
                    DNS:localhost
            :
            :

The signed client certificate is valid for **365 days** and will be used by the
client to authenticate itself to servers that trust the CA. The signing
process proves that the CA has verified the client and endorses the client’s
certificate.

---

## Usage

Clone the repo;

```bash
git clone https://github.com/bilusteknoloji/reverse-ssl-pinning-demo
cd reverse-ssl-pinning-demo/
```

## Run Server

```bash
go run cmd/server/main.go       # runs with default cert paths and port (8443)
go run cmd/server/main.go -h    # see command-line args

  -client-ca string
    	path to the client CA certificate (default "certs/client/client-ca.pem")
  -port string
    	port for the server to listen on (default "8443")
  -server-cert string
    	path to the server certificate (default "certs/server/server-cert.pem")
  -server-key string
    	path to the server private key (default "certs/server/server-key.pem")

```

---

## Client Request with `go`

Now, open another terminal tab and `cd` to cloned repo location.

```bash
cd reverse-ssl-pinning-demo/
go run cmd/client/main.go       # make request to local server with defaults
go run cmd/client/main.go -h    # see command-line args

  -client-ca string
    	path to the client CA certificate (default "certs/client/client-ca.pem")
  -client-cert-signed string
    	path to the signed client certificate (default "certs/client/client-cert-signed.pem")
  -client-key string
    	path to the client private key (default "certs/client/client-key.pem")
  -port string
    	port for the server to connect (default "8443")

```

---

## Client Request with `curl`

Now, open another terminal tab and `cd` to cloned repo location.

```bash
cd reverse-ssl-pinning-demo/
curl --cacert certs/server/server-cert.pem \
     --cert certs/client/client-cert-signed.pem \
     --key certs/client/client-key.pem \
     https://localhost:8443
```

---

## Extras; `python` Client Request

Open another terminal tab and `cd` to cloned repo location, check your `python` version:

```bash
cd reverse-ssl-pinning-demo/
python --version
# mine is Python 3.12.2
```

Now, create your virtual environment:

```bash
python -m venv venv        # create virtual environment
source venv/bin/activate   # activate virtual environment
pip install requests       # install python http client
python client.py           # run
```

Python client is properly configured to trust the **CA** certificate that
signed the server certificate.

---

## Rake Tasks

```bash
$ rake -T

rake release[revision]  # release new version major,minor,patch, default: patch
rake run:curl_client    # run curl client
rake run:go_client      # run go client
rake run:python_client  # run python client
rake run:server         # run server
```

---

## License

This project is licensed under MIT

---


This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].

[coc]: https://github.com/bilusteknoloji/reverse-ssl-pinning-demo/blob/main/CODE_OF_CONDUCT.md
