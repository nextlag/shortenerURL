// Package main provides the command-line tool to shorten URLs using the shortenerURL service.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nextlag/shortenerURL/internal/config"
)

// Setup initializes the configuration, reads a long URL from the console input,
// sends a POST request to the shortenerURL service, and prints the response.
func Setup() {
	if err := config.MakeConfig(); err != nil {
		log.Fatal(err)
	}
	flag.Parse() // Parsing command-line flags
	endpoint := config.Cfg.BaseURL

	// Data container for the request
	data := url.Values{}
	fmt.Println("Введите длинный URL")

	// Open a stream for reading from the console
	reader := bufio.NewReader(os.Stdin)

	// Read a line from the console
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")

	// Fill the data container with the long URL
	data.Set("service", long)

	// Create an HTTP client
	client := &http.Client{}

	// Create a POST request with the long URL in the body
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	// Set the content type header
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request and get the response
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Print the response status code
	fmt.Println("Статус-код ", response.Status)

	// Read and print the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

// main is the entry point of the application.
func main() {
	Setup()
}
