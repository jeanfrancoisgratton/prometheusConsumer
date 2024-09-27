package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Send the POST command to the listener
func sendFileCommand(cfg Config_s, command string) *cerr.CustomError {
	// Initialize PrometheusTarget_s
	pt := PrometheusTarget_s{}

	// If the command is "add" or "rm", set the target hostname
	if command == "add" || command == "rm" {
		hostname, err := os.Hostname()
		// On MacOS, hostnames are suffixed with ".local". We need to trim that
		hostname = strings.TrimSuffix(hostname, ".local")
		if err != nil {
			return &cerr.CustomError{Title: "Error getting the client name",
				Message: fmt.Sprintf("Unable to get hostname: %v\n", err)}
		}
		pt = PrometheusTarget_s{Targets: []string{hostname}}
	}

	// Create the CommandPayload_s with the command and PrometheusTarget
	commandPayload := CommandPayload_s{
		Command:          command,
		PrometheusTarget: pt, // Empty for "ls", filled for "add" and "rm"
	}

	// Marshal the command payload to JSON
	jsonPayload, err := json.Marshal(commandPayload)
	if err != nil {
		return &cerr.CustomError{Title: "Error marshalling payload", Message: err.Error()}
	}

	// Read the CA certificate
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

	// Handle response from "ls" command if necessary
	if command == "ls" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return &cerr.CustomError{Title: "Error reading response body", Message: err.Error()}
		}
		fmt.Println("Response from 'ls' command:", string(body))
	}

	return nil
}
