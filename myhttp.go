package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func main() {
	threadLimit := flag.Int("parallel", 10, "max thread count for connecting")
	flag.Parse()

	hosts := flag.Args()

	// Channel for limit threads
	limitChan := make(chan struct{}, *threadLimit)
	defer func() {
		close(limitChan)
	}()
	// Channel for getting result string for each host
	resultChan := make(chan string, len(hosts))

	var waitGetResponse sync.WaitGroup
	var waitResult sync.WaitGroup

	// Start goroutine for getting results from resultChan
	waitResult.Add(1)
	go getResult(resultChan, &waitResult)

	// In loop we start getting response for each host
	// limitChan block starting new goroutines if maximum number of goroutines started
	// new goroutine started right after one of working goroutine finished
	for _, host := range hosts {
		limitChan <- struct{}{}
		waitGetResponse.Add(1)
		go processHost(host, limitChan, resultChan, &waitGetResponse)
	}
	waitGetResponse.Wait()
	close(resultChan)
	waitResult.Wait()

}

// processHost get response from host and add modified result to channel
func processHost(host string, limitChan chan struct{}, resultChan chan string, waitGroup *sync.WaitGroup) {
	defer func() {
		<-limitChan
		waitGroup.Done()
	}()
	_, err := url.ParseRequestURI(host)
	if err != nil {
		host = "http://" + host

		_, err = url.ParseRequestURI(host)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

	}

	resp, err := http.Get(host)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	resultChan <- fmt.Sprintf("%s %x", host, md5.Sum(body))
}

// getResult read from resultChan and print all getted strings until the channel is closed
func getResult(resultChan chan string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for {
		val, ok := <-resultChan
		if !ok {
			return
		}
		fmt.Println(val)
	}
}
