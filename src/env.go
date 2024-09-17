package main

import (
	"encoding/json"
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"os"
	"path/filepath"
	"strings"
)

type Config_s struct {
	CAcert string `json:"cacert"`
	//	Cert        string `json:"cert"`
	//	Key         string `json:"key"`
	ListenerURL string `json:"listenerurl"`
}

func loadConfig() (Config_s, *cerr.CustomError) {
	var payload Config_s

	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "prometheusConsumer.json")
	_, err := os.Stat(rcFile)
	// We need to create the environment file if it does not exist
	if os.IsNotExist(err) {
		panic("Configuration file not found")
	}

	jFile, err := os.ReadFile(rcFile)
	if err != nil {
		return Config_s{}, &cerr.CustomError{Title: "Error reading the file", Message: err.Error()}
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return Config_s{}, &cerr.CustomError{Title: "Error unmarshalling JSON", Message: err.Error()}
	} else {
		return payload, nil
	}
}

func (cs Config_s) SaveEnvironmentFile() *cerr.CustomError {
	jStream, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return &cerr.CustomError{Title: err.Error(), Fatality: cerr.Fatal}
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "prometheusConsumer.json")
	if err = os.WriteFile(rcFile, jStream, 0644); err != nil {
		return &cerr.CustomError{Title: "Unable to write JSON file", Message: err.Error(), Fatality: cerr.Fatal}
	}
	return nil
}

func setup() *cerr.CustomError {
	cfg := Config_s{}

	cfg.CAcert = hf.GetStringValFromPrompt("Enter the path to your CA certificate: ")
	//	cfg.Cert = hf.GetStringValFromPrompt("Enter the path to your SSL certificate: ")
	//	cfg.Key = hf.GetStringValFromPrompt("Enter the path to its key: ")
	cfg.ListenerURL = hf.GetStringValFromPrompt("Enter the full URL to your listener, including protocol and port: ")

	if !strings.HasPrefix(cfg.ListenerURL, "https://") {
		return &cerr.CustomError{Title: "Invalid URL", Message: "URL must start with https://"}
	}
	return cfg.SaveEnvironmentFile()
}
