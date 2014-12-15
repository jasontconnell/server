package server

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"io/ioutil"
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
	BaseContentDir string
}

func (site *Site) AddHandler(pattern string, handleFunc func(http.ResponseWriter, *http.Request)) {
	var h = &Handler{ pattern: pattern, handler: handleFunc}
	site.Handlers = append(site.Handlers, *h)
}

func static(site Site, w http.ResponseWriter, req *http.Request){
	filePath := site.BaseContentDir + req.URL.Path
	
	if content, err := ioutil.ReadFile(filePath); err == nil {
		writeContent := false
		if etagResult := checkETag(content, w, req); etagResult {
			writeContent = true
		} else if dateResult := checkDate(filePath, w, req); etagResult && dateResult {
			writeContent = true
		}

		if (writeContent){
			w.Header().Add("Content-Type", getMimeType(filePath))
			w.WriteHeader(200)
			sendContent(content, w, req)
		}
	}
}

func makeStatic (site Site) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, req *http.Request){
		static(site, w, req)
	}
}


func Start(site Site){
	mux := mux.NewRouter()

	siteStaticHandler := makeStatic(site)
	mux.HandleFunc("/static/css/{filename}", siteStaticHandler)
	mux.HandleFunc("/static/js/{filename}", siteStaticHandler)

	for _, h := range site.Handlers {
		mux.HandleFunc(h.pattern, h.handler)
	}


	server := &http.Server{
		Addr: site.Domain + ":"  + fmt.Sprint(site.Port),
		Handler: mux,
	}


	fmt.Println("Starting server ", server.Addr)
	log.Fatal(server.ListenAndServe())
}