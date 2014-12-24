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
	AsyncHandlers []Handler
	BaseContentDir string
	Template *template.Template
	State AppState
}

func NewSite(domain string, port int, contentDir string) Site {
	site := Site{ Domain: domain, Port: port, BaseContentDir: contentDir }
	site.State = NewAppState()
	site.Template = template.New("Templates")
	site.Template.Funcs(template.FuncMap{ "html": ToHtml, "toKey": ToKey, "formatDate": FormatDate })

	return site
}

func (site *Site) AddHandler(pattern string, handleFunc func(Site, http.ResponseWriter, *http.Request)) {
	var h = &Handler{ pattern: pattern, handler: handleFunc}
	site.Handlers = append(site.Handlers, *h)
}

func (site *Site) AddAsyncHandler(pattern string, handleFunc func(Site, http.ResponseWriter, *http.Request)) {
	var h = &Handler{ pattern: pattern, handler: handleFunc}
	site.AsyncHandlers = append(site.AsyncHandlers, *h)
}

func (site *Site) AddState(key string, value interface{}){
	if site.State != nil {
		site.State[key] = value
	}
}

func (site *Site) GetState(key string) interface{}{
	var ret interface{}

	if site.State != nil {
		ret = site.State[key]
	}
	return ret
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
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func makeStatic (site Site) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, req *http.Request){
		static(site, w, req)
	}
}

func dynamicHandler (site Site, handler func (Site, http.ResponseWriter, *http.Request), contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", contentType)
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
	mux.HandleFunc(`/static/{path:[a-zA-Z0-9\\/\-\._]+}`, siteStaticHandler)

	for _, h := range site.Handlers {
		mux.HandleFunc(h.pattern, dynamicHandler(site, h.handler, "text/html"))
	}

	for _, h := range site.AsyncHandlers {
		mux.HandleFunc(h.pattern, dynamicHandler(site, h.handler, "application/json"))
	}


	server := &http.Server{
		Addr: site.Domain + ":"  + fmt.Sprint(site.Port),
		Handler: mux,
	}

	fmt.Println("Starting server ", server.Addr)
	log.Fatal(server.ListenAndServe())
}