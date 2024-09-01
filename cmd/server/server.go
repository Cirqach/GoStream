package server

import (
	"net/http"
	"strings"

	"github.com/Cirqach/GoStream/cmd/api/handler"
	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	mw "github.com/Cirqach/GoStream/cmd/middleware"
	queuecontroller "github.com/Cirqach/GoStream/cmd/queueController"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	Host     = "http://localhost:8080"
)

// Server struct    struct allow to control all server
type Server struct {
	ip string                           // server protocol://ip:port addres
	r  *chi.Mux                         // router
	b  *broadcast.Engine                // broadcast engine
	d  *database.DatabaseController     // database controller
	q  *queuecontroller.QueueController // queue controller
}

func NewServer(protocol, ip, port string) *Server {
	logger.LogMessage("NewServer",
		"Creating new server")
	return &Server{
		ip: protocol + "://" + ip + ":" + port,
		r:  chi.NewRouter(),
		b:  broadcast.NewEngine(),
		d:  database.NewDatabaseController(),
	}
}

func (s *Server) StartServer() {
	s.r.Use(middleware.Logger)

	s.startServices()
	s.serveStaticFiles()
	s.handleRoutes()

	err := http.ListenAndServe(":"+strings.Split(s.ip, ":")[2], s.r)
	if err != nil {
		logger.Fatal(logger.GetFuncName(0), err.Error())
	}

}

func (s *Server) startServices() {
	s.d.MakeConnection()
	s.q = queuecontroller.NewQueueController(s.d)
	go s.q.StartControlling(s.b)
	go s.b.Hub.RunHub()

}

func (s *Server) handleRoutes() {
	s.r.Get("/ws", handler.WebsocketHandler(s.b.Hub))
	s.r.Get("/", handler.RootHandler(s.ip))
	s.r.Get("/watch", handler.WatchHandler(s.ip))
	s.r.Get("/login", handler.LoginForm)
	s.r.Post("/auth", handler.Login(s.d))
	s.r.Get("/book", handler.BookTimeFormHandler(s.ip))

	s.r.Route("/", func(r chi.Router) {
		r.Use(mw.Auth())
		r.Post("/book", handler.BookTimeHandler(s.q))
	})
}

func (s *Server) serveStaticFiles() {

	logger.LogMessage(logger.GetFuncName(0), "Serving static files")
	s.r.Handle(
		"/video/processed/*",
		http.StripPrefix("/video/processed/",
			http.FileServer(http.Dir("./video/processed/"))))
	s.r.Handle(
		"/web/static/js/*",
		http.StripPrefix("/web/static/js/",
			http.FileServer(http.Dir("./web/static/js/"))))
	s.r.Handle(
		"/web/static/css/*",
		http.StripPrefix("/web/static/css/",
			http.FileServer(http.Dir("./web/static/css/"))))

}
