package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var NumberWorkers int
	var RaitLimitPerSecond int
	var StatInterval time.Duration
	var URL string
	flag.IntVar(&NumberWorkers, "workers", 5, "Number of workers")
	flag.IntVar(&RaitLimitPerSecond, "raitlimits", 1000, "Rait limit per second")
	flag.DurationVar(&StatInterval, "statinterval", time.Second*2, "Number of workers")
	flag.StringVar(&URL, "URL", "http://localhost:8080/getname", "URL to get")
	flag.Parse()

	exitChan := make(chan os.Signal, 1)
	closeChan := make(chan struct{})
	readyChan := make(chan struct{})

	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	wg := sync.WaitGroup{}
	wg.Add(NumberWorkers)
	for i := 0; i < NumberWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			workerGet(i, StatInterval, URL, readyChan, closeChan)
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		tikerDuration := time.Duration(int(time.Second) / RaitLimitPerSecond)
		tiker := time.NewTicker(tikerDuration)
		defer tiker.Stop()
		for {
			select {
			case <-closeChan:
				return
			case <-tiker.C:
				readyChan <- struct{}{}
			}
		}
	}()

	<-exitChan
	close(closeChan)
	wg.Wait()
}

func workerGet(threadNumber int, statInterval time.Duration, url string, readyChan chan struct{}, chanClose chan struct{}) {
	tiker := time.NewTicker(statInterval)
	defer tiker.Stop()
	answer500 := 0
	answerOK := 0
	answerERR := 0
	lastErr := ""

	for {
		select {
		case <-chanClose:
			return
		case <-tiker.C:
			fmt.Printf("Thread %d, statusOK=%d, status500=%d, err=%d\n", threadNumber, answerOK, answer500, answerERR)
			if answerERR != 0 {
				fmt.Printf("Thread %d, last Err=%s\n", lastErr)
			}
			answer500 = 0
			answerOK = 0
			answerERR = 0
		case <-readyChan:
			client:=http.DefaultClient
			client.Get()
			client.

			resp, err := http.Get(url)
			if err != nil {
				answerERR++
				lastErr = err.Error()
				continue
			}
			if resp.StatusCode == http.StatusOK {
				answerOK++
			} else {
				answerERR++
			}
			resp.Body.Close()
		}
	}
}
