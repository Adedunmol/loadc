package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

type Command struct {
	RequestNumber *int
	URL           *string
	CRequests     *int
	FileLocation  *string
}

type ResponseResult struct {
	Successes      int
	Failures       int
	MinRequestTime float64
	MaxRequestTime float64
}

func (command *Command) makeRequestC(responseResult *ResponseResult, wg *sync.WaitGroup, mux *sync.Mutex) error {

	newChan := make(chan string, 10)

	wg.Add(1)
	go func(wg *sync.WaitGroup, c chan string) {
		defer wg.Done()

		for i := 0; i < *command.RequestNumber; i++ {
			c <- *command.URL
		}

		close(c)

	}(wg, newChan)

	if *command.URL == "" && *command.FileLocation == "" {
		return errors.New("no url to make request to")
	}

	wg.Add(*command.CRequests)
	for i := 0; i < *command.CRequests; i++ {
		go worker(newChan, responseResult, wg, mux)
	}

	return nil
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

func (command *Command) makeRequestSeq(responseResult *ResponseResult) error {

	if *command.URL == "" && *command.FileLocation == "" {
		return errors.New("no url to make request to")
	}

	for i := 0; i < *command.RequestNumber; i++ {
		startTime := time.Now()

		response, err := http.Head(*command.URL)

		endTime := time.Now()

		if err != nil {
			return err
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

	return nil
}

func (command *Command) makeRequestFile(responseResult *ResponseResult, wg *sync.WaitGroup, mux *sync.Mutex) {

	fileDes, err := os.Open(*command.FileLocation)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileDes.Close()

	bufferedReader := bufio.NewScanner(fileDes)

	for bufferedReader.Scan() {
		*command.URL = bufferedReader.Text()

		if *command.CRequests != 0 {
			err := command.makeRequestC(responseResult, wg, mux)

			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err := command.makeRequestSeq(responseResult)

			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func main() {

	command := Command{
		RequestNumber: flag.Int("n", 10, "Number of requests"),
		URL:           flag.String("u", "", "URL to make request(s) to"),
		CRequests:     flag.Int("c", 0, "Number of concurrent requests"),
		FileLocation:  flag.String("f", "", "Location to file to read URLs from"),
	}

	responseResult := ResponseResult{
		Successes:      0,
		Failures:       0,
		MinRequestTime: math.Inf(1),
		MaxRequestTime: math.Inf(-1),
	}

	flag.Parse()

	var err error

	var wg sync.WaitGroup

	mux := &sync.Mutex{}

	if *command.FileLocation == "" {

		if *command.CRequests != 0 {

			err = command.makeRequestC(&responseResult, &wg, mux)

			wg.Wait()

		} else {

			err = command.makeRequestSeq(&responseResult)

		}
	} else {
		command.makeRequestFile(&responseResult, &wg, mux)

		wg.Wait()
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if responseResult.Successes == 0 {
		return
	}

	fmt.Println("Results:")
	fmt.Println(" Total Requests  (2xx)..........: ", responseResult.Successes)
	fmt.Println(" Failed Requests (5xx)..........: ", responseResult.Failures)
	fmt.Println()
	fmt.Println("Total request time (min, max, mean)...: ", responseResult.MinRequestTime, responseResult.MaxRequestTime, roundFloat((responseResult.MaxRequestTime+responseResult.MinRequestTime)/2, 2))
}
