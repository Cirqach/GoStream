package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Cirqach/GoStream/cmd/server"
	"github.com/Cirqach/GoStream/internal/env"
)

var protocol = flag.String("protocol", "http", "protocol for web connection")
var ip = flag.String("ip", "localhost", "ip of server")
var port = flag.String("port", "8080", "port for web service")

func main() {
	env.LoadEnv()

	// TODO: remove it after testing
	fmt.Println("DELETE IT AFTER TESTING")
	fmt.Println(os.Getenv("SECRET_KEY"))
	fmt.Println(os.Getenv("DATABASE_NAME"))
	fmt.Println(os.Getenv("DATABASE_USER"))
	fmt.Println(os.Getenv("DATABASE_PASSWORD"))

	flag.Parse()
	server := server.NewServer(*protocol, *ip, *port)
	server.StartServer()

}
