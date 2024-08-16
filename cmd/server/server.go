package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	videoprocessor "github.com/Cirqach/GoStream/cmd/videoProcessor"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

type Server struct {
	router             *mux.Router
	videoProcessor     *videoprocessor.VideoProcessor
	broadcastEngine    *broadcast.BroadcastEngine
	databaseController *database.DatabaseController
	protocol           string
	ip                 string
	port               string
}

type TemlateData struct {
	Host string
}

func NewServer(protocol, ip, port string) *Server {
	logger.LogMessage("NewServer",
		"Creating new server")
	return &Server{
		router:             mux.NewRouter(),
		videoProcessor:     videoprocessor.NewVideoProcessor(),
		broadcastEngine:    broadcast.NewBroadcastEngine(),
		databaseController: database.NewDatabaseController(),
		protocol:           protocol,
		ip:                 ip,
		port:               port,
	}
}

func (s *Server) StartServer() {

	logger.LogMessage(logger.GetFuncName(0), "Starting server")

	s.databaseController.MakeConnection()

	go s.broadcastEngine.Hub.Run()

	logger.LogMessage(logger.GetFuncName(0), "Handling websocket")
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.websocketHandler(w, r)
	})

	logger.LogMessage(logger.GetFuncName(0), "Serving static files")
	s.router.PathPrefix("/video/processed/").
		Handler(http.StripPrefix("/video/processed/",
			http.FileServer(http.Dir("./video/processed/"))))
	s.router.PathPrefix("/web/static/js/").
		Handler(http.StripPrefix("/web/static/js/",
			http.FileServer(http.Dir("./web/static/js/"))))
	s.router.PathPrefix("/web/static/css/").
		Handler(http.StripPrefix("/web/static/css/",
			http.FileServer(http.Dir("./web/static/css/"))))

	logger.LogMessage(logger.GetFuncName(0), "Serving routes: "+s.protocol+s.ip+s.port)
	s.router.HandleFunc("/", RootHandler(s.protocol+s.ip+s.port))
	s.router.HandleFunc("/watch", WatchHandler(s.protocol+s.ip+s.port))
	s.router.HandleFunc("/book", s.BookatimeHandler(s.protocol+s.ip+s.port))
	s.router.HandleFunc("/auth", s.authHandler)

	logger.LogMessage(logger.GetFuncName(0), "Listen on "+s.port)
	go log.Fatal(http.ListenAndServe(s.port, s.router))
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
		logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("%d %s %s%s by %s", http.StatusOK, r.Method, r.Host, r.URL.Path, r.RemoteAddr))
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

// TODO: need to create websocket connection and send result of processing to client
func (s *Server) BookatimeHandler(host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie("userID")
		logger.LogMessage(logger.GetFuncName(0), "userID: "+userID.Value)
		if err != nil {
			// Handle missing cookie
			logger.LogError(logger.GetFuncName(0), err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication required"))
			return // Close connection after informative message
		}
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
			cookie, err := r.Cookie("userID")
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					http.Error(w, "cookie not found", http.StatusBadRequest)
				default:
					log.Println(err)
					http.Error(w, "server error", http.StatusInternalServerError)
				}
				return
			}

			// Echo out the cookie value in the response body.
			w.Write([]byte(cookie.Value))
			file, handler, err := r.FormFile("file")
			defer file.Close()
			if err = s.videoProcessor.SaveVideo(file, handler); err != nil {
				logger.LogError(logger.GetFuncName(0), err.Error())
			}
			logger.LogMessage(logger.GetFuncName(0), "Video saved")
			time := r.FormValue("time")
			date := r.FormValue("date")
			logger.LogMessage(logger.GetFuncName(0), "extracted time and date: "+date+" "+time)
			if err = s.databaseController.AddVideoToQueue(handler.Filename, date+" "+time); err != nil {
				logger.LogError(logger.GetFuncName(0), "Error adding video to queue: "+err.Error())
			}
			// s.broadcastEngine.Hub.SendToClient(
			// 	s.broadcastEngine.Hub.FindClient(userID.Value),
			// 	[]byte(`
			// 	<div class="alert" role="alert" style="text-align: center; background-color: #d4edda; border-color: #c3e6cb">
			// 		Video was added to queue
			// 	</div>
			//
			// 	`),
			// )
		}
	}
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(w, r, nil)
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
			return
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}
	log.Println("Connection established")
	// Rest of your code using userID
	//TODO: uncomment
	// client := &broadcast.Client{Hub: s.broadcastEngine.Hub,
	// 	Conn: conn,
	// 	Send: make(chan []byte, 256),
	// 	Id:   userID.Value}
	// client.Hub.Register <- client
}

// TODO: change http to https
// TODO: add normal authentication
func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	logger.LogMessage(logger.GetFuncName(0), "Handling auth request")
	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:     "userID",
		Value:    "govno",
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})
	w.Write([]byte("cookie set!"))

}
