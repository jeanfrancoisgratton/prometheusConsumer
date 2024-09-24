package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	var err error
	var cfg Config_s
	var ce *cerr.CustomError
	command := "add"

	if err = os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "JFG"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Define command-line flags for the client
	addFlag := flag.Bool("add", false, "Add current hostname")
	rmFlag := flag.Bool("rm", false, "Remove current hostname")
	setupFlag := flag.Bool("setup", false, "Run setup and exit")
	versionFlag := flag.Bool("version", false, "Displays the version info and exits")
	flag.Parse()

	// -version flag
	if *versionFlag {
		fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), hf.White(fmt.Sprintf("2.00.00-%s 2024.09.23", runtime.GOARCH)))
		os.Exit(0)
	}
	// Check if the "-setup" flag is set
	if *setupFlag {
		// Call the setup function and exit
		if ce = setup(); ce != nil {
			fmt.Println(ce.Error())
		} else {
			return
		}
	}

	// both add and rm flags cannot simultaneously be present or absent
	if *addFlag == *rmFlag {
		fmt.Println("Both -add and -rm cannot simultaneously be set or unset")
		os.Exit(1)
	} else {
		if *rmFlag {
			command = "rm"
		}
	}
	
	// Load the configuration file
	if cfg, ce = loadConfig(); ce != nil {
		fmt.Println(ce.Error())
	}

	// Prepare the request based on the command
	if ce := sendFileCommand(cfg, command); ce != nil {
		fmt.Println(ce.Error())
	}
}

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
	resp, err := client.Post(cfg.ListenerURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error sending HTTPS request: %v", err)
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
