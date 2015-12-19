package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Port          uint   `json:"port"`
	Database      string `json:"database"`
	Auth0Secret   string `json:"auth0_secret"`
	Auth0ClientId string `json:"auth0_client_id"`
	Auth0Domain   string `json:"auth0_domain"`
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
