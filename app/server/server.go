package server

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/Cirqach/GoStream/app/server/client"
	"github.com/Cirqach/GoStream/app/server/hub"
	"github.com/Cirqach/GoStream/app/server/webrtc"
)

var addr = flag.String("addr", "localhost:8080", "url:port for web service")

var upgrader  = websocket.Upgrader{}

func StartServer() {
	flag.Parse()
	hub := NewHub()
	log.SetFlags(0)
	go hub.Run()
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/watch", watchHandler)
	r.HandleFunc("/book", bookatimeHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Serving websocket request")
		serveWs(hub, w, r)
	})

	fmt.Printf("\nListen on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}

func rootHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("%d %s %s%s by %s",http.StatusOK,r.Method,r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/index.html"))
	tmpl.Execute(w, nil)
}

func watchHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("%d %s %s%s by %s",http.StatusOK,r.Method,r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/watch.html"))
	tmpl.Execute(w, nil)
}

func bookatimeHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("%d %s %s%s by %s",http.StatusOK,r.Method,r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/bookatime.html"))
	tmpl.Execute(w, nil)
}

