package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"wee-server/blog"
	"wee-server/support"
)

const (
	help       = "-h"
	help2      = "--help"
	createFile = "-f"
	useConfig  = "-c"
	configFile = "configuration.json"
)
/**
main entry of wee-server
 */
func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		startServerFromFile(configFile)
		return
	case 1:
		a := args[0]
		if a == createFile {
			writeConfigurationToFile()
			return
		} else if a == help || a == help2 {
			printHelp()
			return
		}
	case 2:
		if args[0] == useConfig && args[1] != "" {
			startServerFromFile(args[1])
			return
		}
	}
	printHelp()
}

func printHelp() {
	log.Println(" \n Usage:\n -f : to create configuration.json\n -c [file]: to start server with configuration file\n ")
}

func writeConfigurationToFile() {
	configuration := support.DefaultConfiguration

	bytes, e := json.MarshalIndent(configuration, "", " ")
	if e != nil {
		panic(e)
	}
	e = ioutil.WriteFile(configFile, bytes, 0644)
	if e != nil {
		panic(e)
	}
	log.Printf("Write configuration to configuration.json.")
}

func startServerFromFile(file string) {
	log.Printf("Reading configuration from file:[%s] ", file)
	bytes, e := ioutil.ReadFile(file)
	if e != nil {
		panic(e)
	}
	config := new(support.Configuration)
	e = json.Unmarshal(bytes, config)
	if e != nil {
		panic(e)
	}
	log.Printf("Configuration loaded from %s.", file)
	blog.Start(config)

}
