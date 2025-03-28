package main

import (
	"encoding/json"
	"os"
)

var configFilePath string = "config.json"

type Config struct {
	BaseURL       string        `json:"baseURL"`
	DefaultParams DefaultParams `json:"defaultParams"`
	BackendPort   int           `json:"backendPort"`
	OutputDir     string        `json:"outputDir"`
	MongoDbUri    string        `json:"mongodbUri"`
	ApiKeyIndex   int           `json:"apiKeyIndex"`
	ApiKeys       []ApiKey      `json:"apiKeys"`
}

type ApiKey struct {
	Email  string `json:"email"`
	ApiKey string `json:"apiKey"`
}

type DefaultParams struct {
	Regions    string `json:"regions"`
	Markets    string `json:"markets"`
	OddsFormat string `json:"oddsFormat"`
}

func initConfig() {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
}

func (config *Config) getApiKey() string {
	return config.ApiKeys[config.ApiKeyIndex].ApiKey
}
