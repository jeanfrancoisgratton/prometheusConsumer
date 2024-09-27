package main

import (
	"flag"
	"fmt"
	cerr "github.com/jeanfrancoisgratton/customError"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"os"
	"path/filepath"
	"runtime"
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
	addFlag := flag.Bool("a", false, "Add current hostname")
	rmFlag := flag.Bool("r", false, "Remove current hostname")
	printCfgFlag := flag.Bool("printconfig", false, "Show the current config and exit")
	setupFlag := flag.Bool("setup", false, "Run setup and exit")
	versionFlag := flag.Bool("version", false, "Displays the version info and exit")
	flag.Parse()

	// -version flag
	if *versionFlag {
		fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), hf.White(fmt.Sprintf("2.01.00-%s 2024.09.27", runtime.GOARCH)))
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

	// Prints the contents of the config file (if any) and exits
	if *printCfgFlag {
		if err := printConfig(); err != nil {
			fmt.Println(err.Error())
		}
		os.Exit(0)
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

func printConfig() *cerr.CustomError {
	if cfg, err := loadConfig(); err != nil {
		return err
	} else {
		fmt.Println("CA certificate: ", cfg.CAcert)
		fmt.Println("Listener URL:", cfg.ListenerURL)
	}
	return nil
}
