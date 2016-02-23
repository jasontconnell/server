package server

import (
	"fmt"
	"net/http"
	"log"
	//"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
)

type Handler struct {
	pattern string
	handler func(Site, http.ResponseWriter, *http.Request)
}

type SocketHandler struct {
	pattern string
	handler func(Site, *websocket.Conn)
}

type Site struct {
	Domain string
	Port int
	Handlers []Handler
	AsyncHandlers []Handler
	SocketHandlers []SocketHandler
	BaseContentDir string
	Template *template.Template
	State AppState
	ServerState ServerState
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(req *http.Request) bool { 
		return true
	},
}

func NewSite(config Configuration) Site {
	site := Site{ Domain: config.HostName, Port: config.Port, BaseContentDir: config.ContentLocation }
	site.State = NewAppState()
	site.ServerState = NewServerState()
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

func (site *Site) AddWebsocketHandler(pattern string, handleFunc func(Site, *websocket.Conn)){
	var h = &SocketHandler{ pattern: pattern, handler: handleFunc}
	site.SocketHandlers = append(site.SocketHandlers, *h)
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

func (site *Site) AddServerState(key string, value interface{}){
	if _,ok := site.ServerState[key]; !ok {
		site.ServerState[key] = value
	}
}

func (site *Site) GetServerState(key string) interface{} {
	var ret interface{}

	if _,ok := site.ServerState[key]; ok {
		ret = site.ServerState[key]
	}
	return ret
}

func dynamicHandler (site Site, handler func (Site, http.ResponseWriter, *http.Request), contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", contentType)
		handler(site, w, req)
	}
}

func websocketHandler(site Site, handler func(Site, *websocket.Conn)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request){
		conn, err := upgrader.Upgrade(w, req, nil)

		if err != nil {
			panic(err)
		}

		handler(site, conn)
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
	mux := http.NewServeMux()

	for _, h := range site.Handlers {
		mux.HandleFunc(h.pattern, dynamicHandler(site, h.handler, "text/html"))
	}

	for _, h := range site.AsyncHandlers {
		mux.HandleFunc(h.pattern, dynamicHandler(site, h.handler, "application/json"))
	}

	for _, h := range site.SocketHandlers {
		mux.HandleFunc(h.pattern, websocketHandler(site, h.handler))
	}

	server := &http.Server{
		Addr: site.Domain + ":"  + fmt.Sprint(site.Port),
		Handler: mux,
	}

	fmt.Println("Starting server ", server.Addr)
	log.Fatal(server.ListenAndServe())
}