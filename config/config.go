package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigList struct {
	DbName  string `yaml:"name"`
	Driver  string `yaml:"driver"`
	Port    int    `yaml:"port"`
	LogFile string `yaml:"log"`
}

var Config ConfigList

func init() {
	b, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	var c map[string]ConfigList
	yaml.Unmarshal(b, &c)

	Config = ConfigList{
		DbName:  c["db"].DbName,
		Driver:  c["db"].Driver,
		Port:    c["web"].Port,
		LogFile: c["web"].LogFile,
	}
}
