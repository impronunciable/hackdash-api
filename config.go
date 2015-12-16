package main

import (
	"os"
	"encoding/json"
	"log"
)

type Configuration struct {
	Port		uint	`json:"port"`
	Database	string	`json:"database"`
}

func LoadConfig() *Configuration {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	settings := Configuration{}
	if err = jsonParser.Decode(&settings); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	return &settings
}
