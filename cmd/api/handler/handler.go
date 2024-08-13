package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Cirqach/GoStream/cmd/logger"
)

type TemlateData struct {
	Host string
}

func RootHandler(host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%d %s %s%s by %s", http.StatusOK, r.Method, r.Host, r.URL.Path, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("./web/templates/index.html"))
		data := TemlateData{Host: host}
		err := tmpl.Execute(w, data)
		if err != nil {
			logger.LogError(logger.GetFuncName(RootHandler(host)), err.Error())
		}
	}
}

func WatchHandler(host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%d %s %s%s by %s", http.StatusOK, r.Method, r.Host, r.URL.Path, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("./web/templates/watch.html"))
		data := TemlateData{Host: host}
		err := tmpl.Execute(w, data)
		if err != nil {
			logger.LogError(logger.GetFuncName(WatchHandler(host)), err.Error())
		}
	}
}

func BookatimeHandler(host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			log.Printf("%d %s %s%s by %s", http.StatusOK, r.Method, r.Host, r.URL.Path, r.RemoteAddr)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/html")
			tmpl := template.Must(template.ParseFiles("./web/templates/bookatime.html"))
			data := TemlateData{Host: host}
			err := tmpl.Execute(w, data)
			if err != nil {
				logger.LogError(logger.GetFuncName(BookatimeHandler(host)), err.Error())
			}
		case "POST":

		}
	}
}
