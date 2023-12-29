package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

func main() {

	command := Command{
		RequestNumber: flag.Int("n", 10, "Number of requests"),
		URL:           flag.String("u", "", "URL to make request(s) to"),
		CRequests:     flag.Int("c", 0, "Number of concurrent requests"),
	}

	responseResult := ResponseResult{
		Successes:      0,
		Failures:       0,
		MinRequestTime: math.Inf(1),
		MaxRequestTime: math.Inf(-1),
	}

	flag.Parse()

	if *command.CRequests == 0 {
		makeRequestSeq(command)
		return
	}

	var wg sync.WaitGroup

	sitesChan := make(chan string, 10)
	mux := &sync.Mutex{}

	makeRequestC(command, &responseResult, sitesChan, &wg, mux)

	wg.Wait()
	fmt.Println("Successes: ", responseResult.Successes)
	fmt.Println("Failures: ", responseResult.Failures)
	fmt.Println("Total request time (min, max, mean)", responseResult.MinRequestTime, responseResult.MaxRequestTime, (responseResult.MaxRequestTime+responseResult.MinRequestTime)/2)
}

type Command struct {
	RequestNumber *int
	URL           *string
	CRequests     *int
}

type ResponseResult struct {
	Successes      int
	Failures       int
	MinRequestTime float64
	MaxRequestTime float64
}

func makeRequestC(command Command, responseResult *ResponseResult, c chan string, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()

	wg.Add(1)
	go func(wg *sync.WaitGroup, c chan string) {
		defer wg.Done()

		for i := 0; i < *command.RequestNumber+1; i++ {
			c <- *command.URL
		}

		close(c)

	}(wg, c)

	if *command.URL == "" {
		fmt.Println("no url to make request to")
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
		startTime := time.Now()

		response, err := http.Head(site)

		endTime := time.Now()

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

		mux.Lock()
		responseResult.MinRequestTime = math.Min(roundFloat((endTime.Sub(startTime).Seconds()), 2), responseResult.MinRequestTime)
		responseResult.MaxRequestTime = math.Max(roundFloat((endTime.Sub(startTime).Seconds()), 2), responseResult.MaxRequestTime)
		mux.Unlock()
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func makeRequestSeq(command Command) {

	if *command.URL == "" {
		fmt.Println("no url to make request to")
		return
	}

	responseResult := ResponseResult{
		Successes: 0,
		Failures:  0,
	}

	for i := 0; i < *command.RequestNumber; i++ {
		startTime := time.Now()

		response, err := http.Head(*command.URL)

		endTime := time.Now()

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

		responseResult.MinRequestTime = math.Min(roundFloat((endTime.Sub(startTime).Seconds()), 2), responseResult.MinRequestTime)
		responseResult.MaxRequestTime = math.Max(roundFloat((endTime.Sub(startTime).Seconds()), 2), responseResult.MaxRequestTime)
	}

	fmt.Println("Successes: ", responseResult.Successes)
	fmt.Println("Failures: ", responseResult.Failures)
	fmt.Println("Total request time (min, max, mean)", responseResult.MinRequestTime, responseResult.MaxRequestTime, (responseResult.MaxRequestTime+responseResult.MinRequestTime)/2)
}
