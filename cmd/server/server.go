package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Cirqach/GoStream/cmd/server/api/handler"
	"github.com/Cirqach/GoStream/cmd/server/broadcast"
	"github.com/gorilla/mux"
)


type server struct {
	router *mux.Router
	port string
}

func NewServer(addr string) *server{
	return &server{
		router: mux.NewRouter(),
		port: addr,
	}
}


func (s *server)StartServer() {
	b := broadcast.NewBroadcast()
	b

	s.router.PathPrefix("/web/static/js/").Handler(http.StripPrefix("/web/static/js/", http.FileServer(http.Dir("./web/static/js/"))))
	s.router.PathPrefix("/web/static/css/").Handler(http.StripPrefix("/web/static/css/", http.FileServer(http.Dir("./web/static/css/"))))
	s.router.HandleFunc("/", handler.RootHandler)
	s.router.HandleFunc("/watch", handler.WatchHandler)
	s.router.HandleFunc("/book", handler.BookatimeHandler)
	
	
	fmt.Printf("\nListen on %s\n", s.port)
	log.Fatal(http.ListenAndServe(s.port, s.router))
}

