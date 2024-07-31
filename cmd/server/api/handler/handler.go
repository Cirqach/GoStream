package handler

import (
	"log"
	"html/template"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("%d %s %s%s by %s",http.StatusOK,r.Method,r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./web/templates/index.html"))
	tmpl.Execute(w, nil)
}

func WatchHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("%d %s %s%s by %s",http.StatusOK,r.Method,r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./web/templates/watch.html"))
	tmpl.Execute(w, nil)
}

func BookatimeHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("%d %s %s%s by %s",http.StatusOK,r.Method,r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./web/templates/bookatime.html"))
	tmpl.Execute(w, nil)
}

