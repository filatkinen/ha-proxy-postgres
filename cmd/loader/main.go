package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
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

	chanClose := make(chan struct{})
	for i := 0; i < Threads; i++ {
		go threadQuery(connURL, wg, i, chanClose)
	}

	chanExit := make(chan os.Signal, 1)
	signal.Notify(chanExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-chanExit

	close(chanClose)

	wg.Wait()
}

func threadQuery(connURL string, wg *sync.WaitGroup, threadNumber int, chanClose <-chan struct{}) {
	wg.Add(1)
	defer wg.Done()

	var numberQuery int
	var timeQueries time.Duration

	var badConnections int
	var conn *pgx.Conn
	var err error
	connect := func() error {
		conn, err = pgx.Connect(context.Background(), connURL)
		if err != nil {
			fmt.Printf("Thread %d - Unable to connect to database: %v\n", threadNumber, err)
			return err
		}
		err = conn.Ping(context.Background())
		if err != nil {
			fmt.Printf("Thread %d - Unable to ping  database: %v\n", threadNumber, err)
			return err
		}
		fmt.Printf("Thread %d OK connection to DB:\n", threadNumber)
		return nil
	}

	err = connect()
	if err != nil {
		return
	}

	defer conn.Close(context.Background())

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
			err = queryDialogs(conn)
			timetoTake := time.Since(timeStart)
			timeQueries += timetoTake
			numberQuery++
			if err != nil {
				time.Sleep(TimeSleepAttempt)
				_ = connect()
				badConnections++
				if badConnections > ConnectAttempts {
					return
				}
			}
		}
	}
}

func queryDialogs(conn *pgx.Conn) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("recoverd from panic:%v", e)
		}
	}()

	err = conn.Ping(context.Background())
	if err != nil {
		return err
	}

	var numberUsers int
	err = conn.QueryRow(context.Background(), `select count(user_id) from users`).Scan(&numberUsers)
	if err != nil {
		return err
	}
	if numberUsers == 0 {
		return errors.New("0 users")
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	user := r1.Intn(numberUsers)
	var user_id int

	query := `select user_id from users OFFSET $1 limit 1`
	err = conn.QueryRow(context.Background(), query, user).Scan(&user_id)
	if err != nil {
		return err
	}

	query = `select user_id, friend_id, message from messages where user_id=$1`
	rows, err := conn.Query(context.Background(), query, user_id)
	if err != nil {
		return err
	}
	defer rows.Close()
	var uID, fID int
	var message string
	for rows.Next() {
		err = rows.Scan(
			&uID, &fID, &message)
		if err != nil {
			return err
		}
		_ = uID
		_ = fID
		_ = message
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
