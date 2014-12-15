package server

import (
	"net/http"
	"log"
)

type Handler func(http.ResponseWriter, *http.Request)

type Site struct {
	Domain string
	Port int
	Handlers []Handler
	TemplateFiles []string
}

func Start(site Site){
	server := &http.Server{
		Addr: site.Domain + ":"  + string(site.Port),
	}


	log.Fatal(server.ListenAndServe())
}