package main

import (
	"flag"

	"github.com/Cirqach/GoStream/cmd/server"
)

var protocol = flag.String("protocol", "http://", "protocol for web connection")
var ip = flag.String("ip", "localhost", "ip of server")
var port = flag.String("port", ":8080", "port for web service")

func main() {
	flag.Parse()
	server := server.NewServer(*protocol, *ip, *port)
	server.StartServer()

}
