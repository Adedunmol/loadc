package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	makeRequest()
}

type Command struct {
	RequestNumber *int
	URL           *string
}

type ResponseResult struct {
	Successes int
	Failures  int
}

func makeRequest() {

	command := Command{
		RequestNumber: flag.Int("n", 10, "Number of requests"),
		URL:           flag.String("u", "", "URL to make request(s) to"),
	}

	responseResult := ResponseResult{
		Successes: 0,
		Failures:  0,
	}

	flag.Parse()

	if *command.URL == "" {
		fmt.Println("no url to make request to")
		return
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
