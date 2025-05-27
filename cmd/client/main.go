package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	name := flag.String("name", "test", "Subdomain name")
	local := flag.String("local", "http://localhost:3000", "Local target URL")
	server := flag.String("server", "http://localhost:8080", "Tunnel server URL")
	retryDelay := flag.Int("retry", 5, "Seconds to wait before retrying registration")
	flag.Parse()

	registerURL := fmt.Sprintf("%s/register?name=%s&target=%s", *server, *name, *local)

	// Retry registration until successful. Let's see how this goes.
	for {
		resp, err := http.Get(registerURL)
		if err != nil {
			log.Printf("Failed to register tunnel: %v. Retrying in %d seconds...", err, *retryDelay)
			time.Sleep(time.Duration(*retryDelay) * time.Second)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Registration failed with status: %s. Retrying in %d seconds...", resp.Status, *retryDelay)
			time.Sleep(time.Duration(*retryDelay) * time.Second)
			continue
		}
		break
	}

	log.Printf("Tunnel registered. Access via: %s/%s\n", *server, *name)

	parsedURL, err := url.Parse(*local)
	if err != nil {
		log.Fatalf("Invalid local URL: %v", err)
	}
	if parsedURL.Host == "" {
		log.Fatalf("Please provide a valid host:port in the --local flag")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Greetings from %s!\nRequested Path: %s\n", parsedURL.Host, r.URL.Path)
	})

	log.Printf("Local server running at http://%s\n", parsedURL.Host)
	log.Fatal(http.ListenAndServe(parsedURL.Host, nil))
}
