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

func (site Site) AddFunc(name string, f interface{}) {
	if site.Template == nil {
		panic("Must set template before adding funcs")
	}
	m := template.FuncMap{}
	m[name] = f
	site.Template.Funcs(m)
}

func Start(site Site){
	mux := mux.NewRouter()

	siteStaticHandler := makeStatic(site)
	mux.HandleFunc(`/static/{path:[a-zA-Z0-9\\/\-\.]+}`, siteStaticHandler)

	for _, h := range site.Handlers {
		mux.HandleFunc(h.pattern, dynamicHandler(site, h.handler))
	}

	if site.Template != nil {
		site.Template.Funcs(template.FuncMap{ "html": ToHtml })
	}

	server := &http.Server{
		Addr: site.Domain + ":"  + fmt.Sprint(site.Port),
		Handler: mux,
	}


	fmt.Println("Starting server ", server.Addr)
	log.Fatal(server.ListenAndServe())
}