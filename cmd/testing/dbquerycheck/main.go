package main

import (
	"github.com/filatkinen/ha-proxy-postgres/internal/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load("serv.env")
	if err != nil {
		log.Fatal(err)
	}
	connURL := os.Getenv("CONNURL")
	if connURL == "" {
		log.Fatal("env CONNURL is not set")
	}

	db, err := storage.New(connURL)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	chanExit := make(chan os.Signal, 1)
	signal.Notify(chanExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			_, err = db.SimpleQueryReturnRandomUserName()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	<-chanExit
}
