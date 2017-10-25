package server

import (
	"github.com/jasontconnell/conf"
)

type Configuration struct {
	HostName string `json:"hostname"`
	Port int 		`json:"port"`
	ContentLocation string `json:"contentLocation"`
}

func LoadConfig(file string) Configuration {
	config := Configuration{}
	conf.LoadConfig(file, &config)
	return config
}