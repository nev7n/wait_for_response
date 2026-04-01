package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	version = "dev"
)

func main() {
	var (
		// command flags
		url          = flag.String("url", "http://localhost/", "URL to poll")
		responseCode = flag.Int("code", 200, "Response code to wait for")
		timeout      = flag.Int("timeout", 30000, "Timeout before giving up in ms")
		interval     = flag.Int("interval", 1000, "Interval between polling in ms")
		localhost    = flag.String("localhost", "", "Ip address to use for localhost")
		username     = flag.String("username", "", "HTTP Basic Auth: username")
		password     = flag.String("password", "", "HTTP Basic Auth: password")
	)

	flag.Parse()

	var (
		// runtime
		userAgent = fmt.Sprintf(
			"nev7n/wait_for_response/%s (+https://github.com/nev7n/wait_for_response)",
			version)

		startTime        = time.Now()
		timeoutDuration  = time.Duration(*timeout) * time.Millisecond
		intervalDuration = time.Duration(*interval) * time.Millisecond
		basicAuth        = ""
	)

	fmt.Printf("Polling URL `%s` for response code %d for up to %d ms at %d ms intervals\n",
		*url, *responseCode, *timeout, *interval)

	if *localhost != "" && strings.Contains(*url, "localhost") {
		*url = strings.ReplaceAll(*url, "localhost", *localhost)
	}

	// validate both are set
	if username != nil || password != nil {
		if username == nil {
			fmt.Printf("Missing username, password was set")
			os.Exit(1)
		}

		if password == nil {
			fmt.Printf("Missing password, username was set (%s)", *username)
			os.Exit(1)
		}

		basicAuth = base64.StdEncoding.EncodeToString(
			fmt.Appendf(nil, "%s:%s", *username, *password))
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), intervalDuration)
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, *url, nil)
		if err != nil {
			fmt.Printf("Error: %s", err)
			cancel()
			continue
		}
		req.Header.Set("User-Agent", userAgent)

		if basicAuth != "" {
			req.Header.Set("Authorization", "Basic "+basicAuth)
		}

		client := http.Client{}
		res, err := client.Do(req)
		cancel()
		if err == nil && res.StatusCode == *responseCode {
			fmt.Printf("Response header: %v", res)
			os.Exit(0)
		}
		time.Sleep(intervalDuration)
		elapsed := time.Since(startTime)
		if elapsed > timeoutDuration {
			fmt.Printf("Timed out\n")
			os.Exit(1)
		}
	}
}
