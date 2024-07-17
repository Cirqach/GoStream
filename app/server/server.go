package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)
func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/",rootHandler)
	r.HandleFunc("/watch", watchHandler)
	r.HandleFunc("/bookatime", bookatimeHandler)
	log.Println("Listen on port 8080")
	http.ListenAndServe(":8080", r)
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

}

func bookatimeHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Handle book request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./static/bookatime.html"))
	tmpl.Execute(w, nil)
}

