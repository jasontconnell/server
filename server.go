package server

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

type Handler struct {
	pattern string
	handler func(http.ResponseWriter, *http.Request)
}

type Site struct {
	Domain string
	Port int
	Handlers []Handler
	TemplateFiles []string
}

func (site *Site) AddHandler(pattern string, handleFunc func(http.ResponseWriter, *http.Request)) {
	var h = &Handler{ pattern: pattern, handler: handleFunc}
	site.Handlers = append(site.Handlers, *h)
}


func Start(site Site){
	mux := mux.NewRouter()

	for _, h := range site.Handlers {
		mux.HandleFunc(h.pattern, h.handler)
	}

	server := &http.Server{
		Addr: site.Domain + ":"  + fmt.Sprint(site.Port),
		Handler: mux,
	}


	fmt.Println("Starting server ", server.Addr)
	log.Fatal(server.ListenAndServe())
	fmt.Println("Server started")
}