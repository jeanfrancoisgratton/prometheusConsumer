package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	var cfg Config_s
	var ce *cerr.CustomError

	if err = os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "JFG"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Define command-line flags for the client
	command := flag.String("c", "", "Command to send to the server (add, rm)")
	setupFlag := flag.Bool("setup", false, "Run setup and exit")
	flag.Parse()

	// Check if the "-setup" flag is set
	if *setupFlag {
		// Call the setup function and exit
		if ce = setup(); ce != nil {
			ce.Error()
		} else {
			return
		}
	}

	// Ensure required flags are provided
	if *command == "" {
		fmt.Println("Command is required")
		os.Exit(1)
	}

	// Get the client hostname
	hostname, err := os.Hostname()
	// On MacOS, hostnames are suffixed with ".local". We need to trim that
	hostname = strings.TrimSuffix(hostname, ".local")
	if err != nil {
		log.Fatalf("Unable to get hostname: %v\n", err)
	}

	// Load the configuration file
	if cfg, ce = loadConfig(); ce != nil {
		ce.Error()
	}
	// Load the CA certificate
	caCert, err := os.ReadFile(cfg.CAcert)
	if err != nil {
		log.Fatalf("Unable to load CA cert file: %v\n", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the HTTPS client with the custom CA
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Prepare the request based on the command
	switch strings.ToLower(*command) {
	case "add", "rm":
		sendFileCommand(client, cfg.ListenerURL, *command, hostname)
	default:
		log.Fatalf("Unknown command: %s\n", *command)
	}
}

func sendFileCommand(client *http.Client, serverAddr, command, hostname string) {
	// Construct the request URL and body
	url := fmt.Sprintf("%s/file?cmd=%s&hostname=%s", serverAddr, command, hostname)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v\n", err)
	}

	// Send the POST request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the server response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Server Response: %s\n", string(body))
}
