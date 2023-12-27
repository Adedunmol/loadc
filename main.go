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

func makeRequest() {

	command := Command{
		RequestNumber: flag.Int("n", 10, "Number of requests"),
		URL:           flag.String("u", "", "URL to make request(s) to"),
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

		fmt.Println("Response code: ", response.StatusCode)
	}
}
