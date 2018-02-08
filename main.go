package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)

	// we send urls to check on the urls channel,
	// but only get them on the output channel if
	// they are accepting connections
	urls := make(chan string)
	output := make(chan string)

	// Spin up a bunch of workers
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)

		go func() {
			for url := range urls {
				if isListening(url) {
					output <- url
				}

			}

			wg.Done()
		}()
	}

	// start waiting for output straight away
	go func() {
		// print all of the URLs we get back
		for url := range output {
			fmt.Println(url)
		}
	}()

	// Open the input file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		domain := strings.ToLower(sc.Text())

		// submit http and https versions to be checked
		urls <- "http://" + domain
		urls <- "https://" + domain
		//urls <- "http://" + domain + ":8080"
		//urls <- "https://" + domain + ":8443"
		//urls <- "http://" + domain + ":81"
		//urls <- "http://" + domain + ":591"
		//urls <- "http://" + domain + ":8008"
	}

	// once we've sent all the URLs off we can close the
	// input channel. The workers will finish what they're
	// doing and then call 'Done' on the WaitGroup
	close(urls)

	// Wait until all the workers have finished
	wg.Wait()
	close(output)
}

func isListening(url string) bool {
	var tr = &http.Transport{
		MaxIdleConns:    30,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second * 3,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := client.Do(req)
	//fmt.Println(err)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return false
	}
	return true
}
