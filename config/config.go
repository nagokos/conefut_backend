package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigList struct {
	DBName  string `yaml:"db_name"`
	Driver  string `yaml:"driver"`
	Port    int    `yaml:"port"`
	LogFile string `yaml:"log_file"`
}

var Config ConfigList

func init() {
	b, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}
	yaml.Unmarshal(b, &Config)
	fmt.Println(Config.DBName)
}
