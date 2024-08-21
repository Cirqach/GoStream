package server

import (
	"net/http"

	"github.com/Cirqach/GoStream/cmd/api/handler"
	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	queuecontroller "github.com/Cirqach/GoStream/cmd/queueController"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	Host     = "http://localhost:8080"
)

// Server struct    struct allow to control all server
type Server struct {
	router             *mux.Router
	broadcastEngine    *broadcast.Engine
	databaseController *database.DatabaseController
	queueController    *queuecontroller.QueueController
}

func NewServer() *Server {
	logger.LogMessage("NewServer",
		"Creating new server")
	return &Server{
		router:             mux.NewRouter(),
		broadcastEngine:    broadcast.NewEngine(),
		databaseController: database.NewDatabaseController(),
	}
}

func (s *Server) StartServer(protocol, ip, port string) {

	logger.LogMessage(logger.GetFuncName(0), "Starting server")

	s.databaseController.MakeConnection()

	s.queueController = queuecontroller.NewQueueController(s.databaseController)
	go s.queueController.StartControlling(s.broadcastEngine.Chan, s.broadcastEngine)

	go s.broadcastEngine.Hub.RunHub()

	logger.LogMessage(logger.GetFuncName(0), "Handling websocket")
	s.router.HandleFunc("/ws", handler.WebsocketHandler(&upgrader, s.broadcastEngine))
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

	logger.LogMessage(logger.GetFuncName(0), "Serving routes: "+protocol+ip+port)
	host := protocol + ip + port

	s.router.HandleFunc("/book", handler.BookatimeHandler(host, s.databaseController)).Methods("GET", "POST")
	s.router.HandleFunc("/", handler.RootHandler(host)).Methods("GET")
	s.router.HandleFunc("/watch", handler.WatchHandler(host)).Methods("GET")
	s.router.HandleFunc("/login", handler.LoginPageHandler(host)).Methods("GET")
	s.router.HandleFunc("/auth", handler.LoginHandler(host, s.databaseController)).Methods("POST")

	logger.LogMessage(logger.GetFuncName(0), "Listen on "+port)
	err := http.ListenAndServe(port, s.router)
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
	}

}
