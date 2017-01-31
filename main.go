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
	"time"
)

func main() {
	flag.Parse()

	path := flag.Arg(0)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		domain := strings.ToLower(sc.Text())

		httpURL := "http://" + domain
		httpsURL := "https://" + domain

		if isListening(httpURL) {
			fmt.Println(httpURL)
		}

		if isListening(httpsURL) {
			fmt.Println(httpsURL)
		}

	}

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
		Timeout:       time.Second * 5,
	}
	_, err := client.Get(url)
	if err != nil {
		return false
	}
	return true
}
