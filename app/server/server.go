package server

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "url:port for web service")

var upgrader  = websocket.Upgrader{}

func StartServer() {
	flag.Parse()
	log.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/",rootHandler)
	r.HandleFunc("/watch", watchHandler)
	r.HandleFunc("/bookatime", bookatimeHandler)
	log.Printf("Listen on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}

func rootHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("Host: %s\nBody: %s\n", r.Host, r.Body)
	log.Println("Handle root request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/index.html"))
	tmpl.Execute(w, nil)
}

func watchHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Handle watch request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/watch.html"))
	tmpl.Execute(w, nil)
	c, err : = upgrader.Upgrade(w,r,nil); err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()}
	log.Println("Websocket connection established")
	for{
		
	}
}

func bookatimeHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Handle book request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/bookatime.html"))
	tmpl.Execute(w, nil)
}

