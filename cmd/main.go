package main

import (
	"flag"
	"fmt"
	"log"

	"user-app/pkg/server"
	"user-app/pkg/user/repository"
)

var (
	port         int
	host         string
	clientId     string
	clientSecret string
	templatesDir string

	smtpServerHost string
	smtpServerPort int
	smtpUser       string
	smptPassword   string

	googleApiKey string

	repositoryType string
	mysqlHost      string
	mysqlPort      int
	mysqlDatabase  string
	mysqlUser      string
	mysqlPassword  string
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
	flag.StringVar(&smtpServerHost, "smtpServerHost", "smtp.gmail.com", "smtp server host")
	flag.IntVar(&smtpServerPort, "smtpServerPort", 587, "smtp server port")
	flag.StringVar(&smtpUser, "smtpUser", "glebstepanov1992@gmail.com", "smtp user")
	flag.StringVar(&smptPassword, "smtpPassword", "", "smtp password")
	flag.StringVar(&googleApiKey, "googleApiKey",
		"", "google api key")
	flag.StringVar(&repositoryType, "repositoryType", "mysql", "type of repository memory or mysql")
	flag.StringVar(&mysqlHost, "mysqlHost", "localhost", "host of mysql")
	flag.IntVar(&mysqlPort, "mysqlPort", 3306, "mysql port")
	flag.StringVar(&mysqlUser, "mysqlUser", "root", "mysql user")
	flag.StringVar(&mysqlPassword, "mysqlPassword", "1234", "mysql password")
	flag.StringVar(&mysqlDatabase, "mysqlDatabase", "userapp", "mysql database name")
	flag.Parse()
	server.InitOAuth(clientId, clientSecret)
	addr := fmt.Sprintf("%s:%d", host, port)
	opts := repository.Opts{
		Mysql: repository.MySQLOpts{
			Host:         mysqlHost,
			Port:         mysqlPort,
			User:         mysqlUser,
			Password:     mysqlPassword,
			DatabaseName: mysqlDatabase,
		},
		Memory: repository.MemoryOpts{},
	}
	srv, err := server.NewServer(addr, templatesDir, repositoryType, opts, smtpServerHost,
		smtpServerPort, smtpUser, smptPassword, googleApiKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Run server on %s\n", addr)
	log.Fatal(srv.Run())
}
