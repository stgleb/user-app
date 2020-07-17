package main

import (
	"flag"
	"fmt"
	"log"

	"user-app/pkg/server"
)

var (
	port           int
	host           string
	clientId       string
	clientSecret   string
	templatesDir   string
	smtpServerHost string
	smtpServerPort int
	smtpUser       string
	smptPassword   string
)

func main() {
	flag.IntVar(&port, "port", 8080, "port number")
	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.StringVar(&clientId, "client_id", "",
		"google client_id")
	flag.StringVar(&clientSecret, "client_secret", "",
		"google client_secret")
	flag.StringVar(&templatesDir, "templatesDir", "templates",
		"templates dir path")
	flag.StringVar(&smtpServerHost, "smtpServerHost", "smtp.mailtrap.io", "smtp server host")
	flag.IntVar(&smtpServerPort, "smtpServerPort", 2525, "smtp server port")
	flag.StringVar(&smtpUser, "smtpUser", "967fe121c1f173", "smtp user")
	flag.StringVar(&smptPassword, "smtpPassword", "4ecbcd773762b7", "smtp password")
	flag.Parse()
	server.InitOAuth(clientId, clientSecret)
	addr := fmt.Sprintf("%s:%d", host, port)
	srv, err := server.NewServer(addr, templatesDir, smtpServerHost,
		smtpServerPort, smtpUser, smptPassword)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Run server on %s\n", addr)
	log.Fatal(srv.Run())
}
