package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Cirqach/GoStream/cmd/server/api/handler"
	"github.com/Cirqach/GoStream/cmd/server/broadcast"
	"github.com/gorilla/mux"
)


type server struct {
	router *mux.Router
	port string
}

func NewServer(addr string) *server{
	log.Println("Creating new server")
	return &server{
		router: mux.NewRouter(),
		port: addr,
	}
}


func (s *server)StartServer() {
	log.Println("Starting server")
	b := broadcast.NewBroadcast()
	go b.Hub.Run()
	log.Println("Handling websocket")
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
b.HandleWebsocket(w, r)
time.Sleep(10 * time.Second)
		log.Println("Putting data to stream")
		b.Hub.Stream <- []byte("Hello world")
	})
	log.Println("Serving static files")
	s.router.PathPrefix("/web/static/js/").Handler(http.StripPrefix("/web/static/js/", http.FileServer(http.Dir("./web/static/js/"))))
	s.router.PathPrefix("/web/static/css/").Handler(http.StripPrefix("/web/static/css/", http.FileServer(http.Dir("./web/static/css/"))))
	log.Println("Serving routes files")
	s.router.HandleFunc("/", handler.RootHandler)
	s.router.HandleFunc("/watch", handler.WatchHandler)
	s.router.HandleFunc("/book", handler.BookatimeHandler)
	
	
	fmt.Printf("\nListen on %s\n", s.port)
	log.Fatal(http.ListenAndServe(s.port, s.router))
}

