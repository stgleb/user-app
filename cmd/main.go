package main

import (
	"flag"
)

var (
	port int
	host string
)

func main(){
	flag.IntVar(&port, "port", 8080, "port number")
	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.Parse()
}
