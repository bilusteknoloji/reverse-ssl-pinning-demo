package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultReadHeaderTimeout = 5 * time.Second

func main() {
	clientCAPath := flag.String("client-ca", "certs/client/client-ca.pem", "path to the client CA certificate")
	serverCertPath := flag.String("server-cert", "certs/server/server-cert.pem", "path to the server certificate")
	serverKeyPath := flag.String("server-key", "certs/server/server-key.pem", "path to the server private key")
	port := flag.String("port", "8443", "port for the server to listen on")

	flag.Parse()

	clientCA, err := os.ReadFile(*clientCAPath)
	if err != nil {
		log.Fatalf("unable to read client CA certificate: %v", err)
	}

	clientCAPool := x509.NewCertPool()
	if ok := clientCAPool.AppendCertsFromPEM(clientCA); !ok {
		log.Fatalf("failed to append client CA certificate")
	}

	tlsConfig := &tls.Config{
		ClientCAs:  clientCAPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	server := &http.Server{
		Addr:              ":" + *port,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		TLSConfig:         tlsConfig,
		Handler:           http.HandlerFunc(handler),
	}

	log.Printf("starting server on https://localhost:%s", *port)

	if errr := server.ListenAndServeTLS(*serverCertPath, *serverKeyPath); errr != nil {
		log.Fatalf("failed to start server: %v", errr)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, ssl pinned client! (path: %s)\n", r.URL)
}
