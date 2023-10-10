package main

import (
	"flag"
	"github.com/filatkinen/ha-proxy-postgres/internal/backend"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port binding")

	err := godotenv.Load("serv.env")
	if err != nil {
		log.Fatal(err)
	}
	connURL := os.Getenv("CONNURL")
	if connURL == "" {
		log.Fatal("env CONNURL is not set")
	}

	serv, err := backend.New(port, connURL)
	if err != nil {
		log.Fatalf("got error creating HTTP Server %s", err)
		return
	}
	log.Printf("Successufuly creating HTTP Server %s", serv.GetID())
	defer func() {
		if err := serv.Close(); err != nil {
			log.Printf("git error closing HTTP Server :%s", err)
		}
	}()

	chanErrorStartingServer := make(chan struct{})
	log.Printf("Starting HTTP Server  %s", serv.GetID())
	go func() {
		err := serv.Start()
		if err != nil {
			chanErrorStartingServer <- struct{}{}
		}
	}()

	chanExit := make(chan os.Signal, 1)
	signal.Notify(chanExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case <-chanErrorStartingServer:
	case <-chanExit:
		err = serv.Stop()
		if err != nil {
			log.Printf("got error stopping HTTP Server: %s", err)
		}
	}
}
