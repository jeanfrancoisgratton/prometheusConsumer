package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"log"
	"net/http"
	"os"
	"strings"
)

func sendFileCommand(cfg Config_s, command string) *cerr.CustomError {
	hostname, err := os.Hostname()
	// On MacOS, hostnames are suffixed with ".local". We need to trim that
	hostname = strings.TrimSuffix(hostname, ".local")
	if err != nil {
		return &cerr.CustomError{Title: "Error getting the client name",
			Message: fmt.Sprintf("Unable to get hostname: %v\n", err)}
	}

	// JSON payload
	pt := PrometheusTarget_s{Targets: []string{hostname}}
	commandPayload := CommandPayload_s{Command: command, PrometheusTarget: pt}

	jsonPayload, err := json.Marshal(commandPayload)
	if err != nil {
		return &cerr.CustomError{Title: "Error marshalling payload", Message: err.Error()}
	}

	caCert, err := os.ReadFile(cfg.CAcert)
	if err != nil {
		return &cerr.CustomError{Title: "Error reading CA cert file", Message: err.Error()}
	}

	// Create a CA certificate pool
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return &cerr.CustomError{Title: "Error parsing CA cert file", Message: "Failed to append CA cert"}
	}

	// Configure TLS settings for HTTPS client
	tlsConfig := &tls.Config{
		RootCAs: caCertPool, // Trust only this CA
	}

	// Create HTTPS client with TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Send the JSON payload via HTTPS POST
	resp, perr := client.Post(cfg.ListenerURL, "application/json", bytes.NewBuffer(jsonPayload))
	if perr != nil {
		log.Fatalf("Error sending HTTPS request: %v", perr)
	}
	defer resp.Body.Close()

	// Log response status
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Payload sent successfully!")
	} else {
		return &cerr.CustomError{Title: "Failed to send payload", Message: resp.Status}
	}
	return nil
}
