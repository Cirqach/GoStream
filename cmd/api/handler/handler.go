package handler

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/auth"
	"github.com/Cirqach/GoStream/internal/database"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
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
			logger.LogError(logger.GetFuncName(0), err.Error())
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
			logger.LogError(logger.GetFuncName(0), err.Error())
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
				logger.LogError(logger.GetFuncName(0), err.Error())
			}
		case "POST":

		}
	}
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	log.Printf("%d %s %s%s by %s", http.StatusOK, r.Method, r.Host, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./web/templates/login.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
	}
}
func Login(d *database.DatabaseController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")
		if d.VerifyUser(username, password) {
			token := auth.GenerateToken()
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    token,
				Expires:  time.Now().Add(time.Hour * 72),
				HttpOnly: true,
			})
		}
	}
}

func WebsocketHandler(hub *broadcast.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		client := &broadcast.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
		client.Hub.Register <- client

		go client.WritePump()
	}
}
