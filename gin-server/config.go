package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`
	JWT struct {
		Secret string `json:"secret"`
	} `json:"jwt"`
}

var AppConfig *Config

func LoadConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &Config{}
	err = decoder.Decode(cfg)

	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	AppConfig = cfg
	log.Println("Config loaded successfully")
}
