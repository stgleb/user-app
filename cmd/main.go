package main

import (
	"flag"
	"fmt"
	"log"
	"user-app/pkg/server"
)

var (
	port int
	host string
	clientId string
	clientSecret string
	templatesDir string
)

func main(){
	flag.IntVar(&port, "port", 8080, "port number")
	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.StringVar(&clientId, "client_id", "436250024602-ia4g4uq0uj14t21snoquq0dr094ivi7k.apps.googleusercontent.com",
		"google client_id")
	flag.StringVar(&clientSecret, "client_secret", "7Oso-KMGSbvt0ksiAHuJAbCS",
		"google client_secret")
	flag.StringVar(&templatesDir, "templatesDir", "templates",
		"templates dir path")
	flag.Parse()
	server.InitOAuth(clientId, clientSecret)
	addr := fmt.Sprintf("%s:%d", host, port)
	srv, err := server.NewServer(addr, templatesDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Run server on %s\n", addr)
	log.Fatal(srv.Run())
}
