package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultTimeout = 5 * time.Second

func main() {
	clientCAPath := flag.String("client-ca", "certs/client/client-ca.pem", "path to the client CA certificate")
	clientCertPath := flag.String(
		"client-cert-signed",
		"certs/client/client-cert-signed.pem",
		"path to the signed client certificate",
	)
	clientKeyPath := flag.String("client-key", "certs/client/client-key.pem", "path to the client private key")
	port := flag.String("port", "8443", "port for the server to connect")

	flag.Parse()

	serverCert, err := os.ReadFile(*clientCAPath)
	if err != nil {
		log.Fatalf("unable to read server certificate: %v", err)
	}

	serverCertPool := x509.NewCertPool()
	if ok := serverCertPool.AppendCertsFromPEM(serverCert); !ok {
		log.Fatalf("failed to append server certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(*clientCertPath, *clientKeyPath)
	if err != nil {
		log.Fatalf("unable to load client certificate/key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      serverCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	body, err := makeRequest(client, port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server response:\n%s\n", body)
}

func makeRequest(c *http.Client, p *string) ([]uint8, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://localhost:"+*p, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer func() {
		if errr := resp.Body.Close(); errr != nil {
			log.Printf("body close err: %v", errr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return body, nil
}
