package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	sitesChan := make(chan string, 10)
	mux := &sync.Mutex{}

	command := Command{
		RequestNumber: flag.Int("n", 10, "Number of requests"),
		URL:           flag.String("u", "", "URL to make request(s) to"),
		CRequests:     flag.Int("c", 0, "Number of concurrent requests"),
	}

	responseResult := ResponseResult{
		Successes: 0,
		Failures:  0,
	}

	flag.Parse()

	wg.Add(1)
	go func(wg *sync.WaitGroup, c chan string) {
		defer wg.Done()

		for i := 0; i < *command.RequestNumber+1; i++ {
			sitesChan <- *command.URL
		}

		close(c)

	}(&wg, sitesChan)

	makeRequest(command, &responseResult, sitesChan, &wg, mux)

	wg.Wait()
	fmt.Println("Successes: ", responseResult.Successes)
	fmt.Println("Failures: ", responseResult.Failures)
}

type Command struct {
	RequestNumber *int
	URL           *string
	CRequests     *int
}

type ResponseResult struct {
	Successes int
	Failures  int
}

func makeRequest(command Command, responseResult *ResponseResult, c <-chan string, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	if *command.URL == "" {
		fmt.Println("no url to make request to")
		return
	}

	if *command.CRequests == 0 {
		makeRequestSeq(command)
		return
	}

	for i := 0; i < *command.CRequests; i++ {
		wg.Add(1)
		go worker(c, responseResult, wg, mux)
	}

	return
}

func worker(sitesChan <-chan string, responseResult *ResponseResult, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	for site := range sitesChan {
		response, err := http.Head(site)

		if err != nil {
			fmt.Println(err)
			return
		}

		if response.StatusCode < 299 {
			mux.Lock()
			responseResult.Successes += 1
			mux.Unlock()
		}

		if response.StatusCode >= 500 {
			mux.Lock()
			responseResult.Failures += 1
			mux.Unlock()
		}
	}
}

func makeRequestSeq(command Command) {

	responseResult := ResponseResult{
		Successes: 0,
		Failures:  0,
	}

	for i := 0; i < *command.RequestNumber; i++ {
		response, err := http.Head(*command.URL)

		if err != nil {
			fmt.Println(err)
			return
		}

		if response.StatusCode < 299 {
			responseResult.Successes += 1
		}

		if response.StatusCode >= 500 {
			responseResult.Failures += 1
		}
	}

	fmt.Println("Successes: ", responseResult.Successes)
	fmt.Println("Failures: ", responseResult.Failures)
}
