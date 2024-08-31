package handler

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	queuecontroller "github.com/Cirqach/GoStream/cmd/queueController"
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

func BookTimeFormHandler(host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
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

func BookTimeHandler(q *queuecontroller.QueueController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		time := r.FormValue("time")
		date := r.FormValue("date")
		if time == "" && date == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("time and date are required"))
			logger.LogError(logger.GetFuncName(0), "time and date are empty")
			return
		}
		// Access the uploaded file
		file, handler, err := r.FormFile("videofile") // Replace "file" with the input field name
		if err != nil {
			http.Error(w, "No file uploaded", http.StatusBadRequest)
			logger.LogError(logger.GetFuncName(0), err.Error())
			return
		}
		defer file.Close()

		var supportedFormats = []string{"mp4", "mkv", "avi", "mov", "wmv", "flv", "webm"}
		extension := strings.Split(handler.Filename, ".")[1]

		if !func(supportedFormats []string, extension string) bool {
			for _, format := range supportedFormats {
				if format == extension {
					return true
				}
			}
			return false
		}(supportedFormats, extension) {
			logger.LogError(logger.GetFuncName(0), "Unsupported file format: "+handler.Filename)
			http.Error(w, "Unsupported file format", http.StatusBadRequest)
			return
		}

		err = saveFile(file, handler)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.LogError(logger.GetFuncName(0), err.Error())
			return
		}

		err = q.BookATime(time, date, handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.LogError(logger.GetFuncName(0), err.Error())
			return
		}

	}
}

func saveFile(file multipart.File, handler *multipart.FileHeader) error {
	// Create a new file on the server
	newFileName := filepath.Join("./video/unprocessed/", handler.Filename) // Adjust the upload directory
	newFile, err := os.Create(newFileName)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		return fmt.Errorf("Error copying file: %v", err)
	}
	return nil
}
func LoginForm(w http.ResponseWriter, r *http.Request) {
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
