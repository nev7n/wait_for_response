package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	var url = flag.String("url", "http://localhost/", "URL to poll")
	var responseCode = flag.Int("code", 200, "Response code to wait for")
	var timeout = flag.Int("timeout", 2000, "Timeout before giving up in ms")
	var interval = flag.Int("interval", 200, "Interval between polling in ms")
	flag.Parse()

	fmt.Printf("Polling URL `%s` for response code %d for up to %d ms at %d ms intervals\n", *url, *responseCode, *timeout, *interval)
	startTime := time.Now()
	timeoutDuration := time.Duration(*timeout) * time.Millisecond
	sleepDuration := time.Duration(*interval) * time.Millisecond
	for {
		res, err := http.Head(*url)
		if err == nil && res.StatusCode == *responseCode {
			fmt.Printf("Response header: %v", res)
			os.Exit(0)
		}
		time.Sleep(sleepDuration)
		elapsed := time.Now().Sub(startTime)
		if elapsed > timeoutDuration {
			fmt.Printf("Timed out\n")
			os.Exit(1)
		}
	}
}
