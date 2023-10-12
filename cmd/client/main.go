package main

import (
	"flag"
	"fmt"
	"io"
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
	flag.IntVar(&RaitLimitPerSecond, "raitlimits", 10000, "Rait limit per second")
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

	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := defaultTransportPointer.Clone() // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100
	//defaultTransport.MaxConnsPerHost = 100
	//defaultTransport.IdleConnTimeout = 1

	myClient := &http.Client{Transport: defaultTransport}

	for {
		select {
		case <-chanClose:
			return
		case <-tiker.C:
			fmt.Printf("Thread %d, statusOK=%d, status500=%d, err=%d\n", threadNumber, answerOK, answer500, answerERR)
			if answerERR != 0 {
				fmt.Printf("Thread %d, last Err=%s\n", threadNumber, lastErr)
			}
			answer500 = 0
			answerOK = 0
			answerERR = 0
			lastErr = ""
		case <-readyChan:
			resp, err := myClient.Get(url)
			if err != nil {
				answerERR++
				lastErr = err.Error()
				continue
			}
			_, _ = io.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusOK {
				answerOK++
			} else {
				answer500++
			}
			resp.Body.Close()
		}
	}
}
