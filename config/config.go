package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigList struct {
	DBName       string `yaml:"name"`
	Driver       string `yaml:"driver"`
	Port         int    `yaml:"port"`
	LogFile      string `yaml:"log"`
	LogErrorFile string `yaml:"errorLog"`
}

var Config ConfigList

func init() {
	b, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	var c map[string]ConfigList
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Printf("Failed to file unmarshal: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		DBName:       c["db"].DBName,
		Driver:       c["db"].Driver,
		Port:         c["web"].Port,
		LogFile:      c["web"].LogFile,
		LogErrorFile: c["web"].LogErrorFile,
	}
}
