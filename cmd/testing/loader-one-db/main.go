package main

import (
	"github.com/filatkinen/ha-proxy-postgres/internal/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	Threads          = 20
	TimeSleepQuery   = time.Millisecond * 1
	TimeStat         = time.Second * 10
	ConnectAttempts  = 100
	TimeSleepAttempt = time.Second * 3
)

func main() {
	err := godotenv.Load("loader.env")
	if err != nil {
		log.Fatal(err)
	}
	connURL := os.Getenv("CONNURL")
	if connURL == "" {
		log.Printf("env CONNURL is not set")
		return
	}

	wg := &sync.WaitGroup{}

	db, err := storage.New(connURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cnanTheadReturn := make(chan int)
	chanClose := make(chan struct{})
	for i := 0; i < Threads; i++ {
		go threadQuery(db, wg, i, chanClose, cnanTheadReturn)
	}

	chanExit := make(chan os.Signal, 1)
	signal.Notify(chanExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	exitCondition := false
	for !exitCondition {
		select {
		case <-chanExit:
			exitCondition = true
		case n := <-cnanTheadReturn:
			log.Printf("recoverd Thred %d", n)
			go threadQuery(db, wg, n, chanClose, cnanTheadReturn)
		}
	}

	close(chanClose)
	wg.Wait()
}

func threadQuery(db *storage.Database, wg *sync.WaitGroup, threadNumber int,
	chanClose <-chan struct{}, chanTread chan<- int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Tread %d recoverd from panic:%s", threadNumber, err)
			chanTread <- threadNumber
		}
	}()

	wg.Add(1)
	defer wg.Done()

	var numberQuery int
	var timeQueries time.Duration

	var badConnections int
	var err error
	if err != nil {
		return
	}

	ticker := time.NewTicker(TimeSleepQuery)
	defer ticker.Stop()

	tickerStat := time.NewTicker(TimeStat)
	defer tickerStat.Stop()

	for {
		select {
		case <-chanClose:
			return
		case <-tickerStat.C:
			log.Printf("Stat: Thread %d, numberQueries=%d, averageTime=%v errors: %d",
				threadNumber, numberQuery, time.Duration(int64(timeQueries)/int64(numberQuery)), badConnections)
			numberQuery = 0
			timeQueries = 0
		case <-ticker.C:
			timeStart := time.Now()
			_, err = db.SimpleQueryReturnRandomUserName()
			timetoTake := time.Since(timeStart)
			timeQueries += timetoTake
			numberQuery++
			if err != nil {
				time.Sleep(TimeSleepAttempt)
				log.Printf("Error thread %d :%s", threadNumber, err)
				badConnections++
			}
		}
	}
}
