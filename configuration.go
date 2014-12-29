package server

import (
	"conf"
)

type Configuration struct {
	HostName string `json:"hostname"`
	Port int 		`json:"port"`
	ContentLocation string `json:"contentLocation"`
	Aliases []string `json:"aliases"`
}

func LoadConfig(file string) Configuration {
	config := Configuration{}
	conf.LoadConfig(file, &config)
	return config
}