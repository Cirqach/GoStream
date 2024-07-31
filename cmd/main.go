package main

import (
	"flag"

	"github.com/Cirqach/GoStream/cmd/server"
)

var addr = flag.String("addr", "localhost:8080", "url:port for web service")

func main() {
	flag.Parse()
	server := server.NewServer(*addr)
	server.StartServer()

}