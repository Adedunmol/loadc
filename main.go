package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	makeRequest()
}

func makeRequest() {
	rawURL := os.Args[1]

	response, err := http.Head(rawURL)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Response code: ", response.StatusCode)
}
