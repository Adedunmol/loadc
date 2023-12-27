package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	makeRequest()
}

func makeRequest() {
	noOfRequests := flag.Int("n", 10, "Number of requests")
	url := flag.String("u", "", "URL to make request(s) to")

	flag.Parse()

	if *url == "" {
		fmt.Println("no url to make request to")
		return
	}

	for i := 0; i < *noOfRequests; i++ {
		response, err := http.Head(*url)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Response code: ", response.StatusCode)
	}
}
