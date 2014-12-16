package server

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"html/template"
)

type Handler struct {
	pattern string
	handler func(Site, http.ResponseWriter, *http.Request)
}

type Site struct {
	Domain string
	Port int
	Handlers []Handler
	BaseContentDir string
	Template *template.Template
}

func (site *Site) AddHandler(pattern string, handleFunc func(Site, http.ResponseWriter, *http.Request)) {
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
			sendContent(content, w, req)
		}
	}
}

func makeStatic (site Site) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, req *http.Request){
		static(site, w, req)
	}
}

func dynamicHandler (site Site, handler func (Site, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request){
		handler(site, w, req)
	}
}


func Start(site Site){
	mux := mux.NewRouter()

	siteStaticHandler := makeStatic(site)
	mux.HandleFunc(`/static/{path:[a-zA-Z0-9\\/\-\.]+}`, siteStaticHandler)
	//mux.HandleFunc("/static/js/{filename}", siteStaticHandler)

	for _, h := range site.Handlers {
		mux.HandleFunc(h.pattern, dynamicHandler(site, h.handler))
	}

	if site.Template != nil {
		funcMap := template.FuncMap {
			"html" : ToHtml,
		}
		site.Template.Funcs(funcMap)
	}

	server := &http.Server{
		Addr: site.Domain + ":"  + fmt.Sprint(site.Port),
		Handler: mux,
	}


	fmt.Println("Starting server ", server.Addr)
	log.Fatal(server.ListenAndServe())
}