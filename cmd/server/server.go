package server

import (
	"net/http"

	"github.com/Cirqach/GoStream/cmd/api/handler"
	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/cmd/middleware"
	videoprocessor "github.com/Cirqach/GoStream/cmd/videoProcessor"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	Host     = "http://localhost:8080"
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
	s.router.HandleFunc("/ws", handler.WebsocketHandler(&upgrader))
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
	host := s.protocol + s.ip + s.port

	s.router.HandleFunc("/book",
		middleware.AuthMiddleware(
			handler.BookatimeHandler(host,
				s.videoProcessor,
				s.databaseController)))
	s.router.HandleFunc("/", handler.RootHandler(host))
	s.router.HandleFunc("/watch", handler.WatchHandler(host))
	s.router.HandleFunc("/cookiet", handler.GetCookie)
	s.router.HandleFunc("/auth", handler.AuthHandler(host))

	logger.LogMessage(logger.GetFuncName(0), "Listen on "+s.port)
	err := http.ListenAndServe(s.port, s.router)
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
	}
}

func hdnler(w http.ResponseWriter, r *http.Request) {}
