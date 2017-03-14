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

func worker(urls, out chan string, wg sync.WaitGroup) {
	defer wg.Done()

	for {
		url, ok := <-urls
		if !ok {
			return
		}
		if isListening(url) {
			out <- url
		}
	}
}

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
		go worker(urls, output, wg)
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       time.Second * 1,
	}
	_, err := client.Get(url)
	if err != nil {
		return false
	}
	return true
}
