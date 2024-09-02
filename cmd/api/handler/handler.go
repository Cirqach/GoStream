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
	"github.com/Cirqach/GoStream/cmd/videoProcessor/ffmpeg"
	"github.com/Cirqach/GoStream/internal/auth"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/google/uuid"

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

// BookTimeHandler function    handler for post request for booking time
func BookTimeHandler(q *queuecontroller.QueueController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// getting values from request form
		time := r.FormValue("time")
		date := r.FormValue("date")
		if time == "" && date == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("time and date are required"))
			logger.LogError(logger.GetFuncName(0), "time and date are empty")
			return
		}
		logger.LogMessage(logger.GetFuncName(0), "extracted time: "+time+" and date: "+date)
		// getting file from request
		file, handler, err := r.FormFile("file") // Replace "file" with the input field name
		if err != nil {
			http.Error(w, "No file uploaded", http.StatusBadRequest)
			logger.LogError(logger.GetFuncName(0), "No file uploaded: "+err.Error())
			return
		}
		logger.LogMessage(logger.GetFuncName(0), "File received: "+handler.Filename)
		defer file.Close()

		// Check if the file format is supported
		logger.LogMessage(logger.GetFuncName(0), "Checking file format")
		var supportedFormats = []string{"mp4", "mkv", "avi", "mov", "wmv", "flv", "webm"}
		extension := strings.Split(handler.Filename, ".")[1]
		if !func(supportedFormats []string, extension string) bool {
			// go in the loop through supported formats
			for _, format := range supportedFormats {
				// if format is supported return true
				if format == extension {
					logger.LogMessage(logger.GetFuncName(0), "File format supported")
					return true
				}
			}
			// format is not supported
			return false
		}(supportedFormats, extension) {
			logger.LogError(logger.GetFuncName(0), "Unsupported file format: "+handler.Filename)
			http.Error(w, "Unsupported file format", http.StatusBadRequest)
			return
		}

		// creating unique name for video file
		filename := uuid.New().String()
		// saving video file to /video/unprocessed directory
		err = saveFile(file, filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.LogError(logger.GetFuncName(0), err.Error())
			return
		}

		videoFilePath := "./video/unprocessed/" + filename
		// getting video duration
		videoDuration, err := ffmpeg.GetVideoDuration(videoFilePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.LogError(logger.GetFuncName(0), "Error getting video duration: "+err.Error())
			return
		}

		// cutting default time.Time format to only HH:mm:ss
		duration := strings.Split(videoDuration.String(), " ")[1]

		// creating new record in database
		logger.LogMessage(logger.GetFuncName(0), "Creating new record in database with data time: "+time+" and date: "+date+" and duration: "+duration)
		err = q.BookATime(time, date, filename, duration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.LogError(logger.GetFuncName(0), "Error booking time: "+err.Error())
			return
		}
		logger.LogMessage(logger.GetFuncName(0), "Time booked successfully")

		// returning success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("time booked"))
		logger.LogMessage(logger.GetFuncName(0), "time booked")

	}
}

func saveFile(file multipart.File, filename string) error {
	logger.LogMessage(logger.GetFuncName(0), "Saving file: "+filename)
	// Create a new file on the server
	newFileName := filepath.Join("./video/unprocessed/", filename) // Adjust the upload directory
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
