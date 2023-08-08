package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

func ScrapeCBREListings() {
	url := "https://www.cbre.ca/properties/office?aspects=isSale"

	resp, err := makeCBRERequest(url)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		fmt.Println(href)
	})

	// Sleep for a while to ensure the program doesn't exit immediately
	time.Sleep(10 * time.Second)
}

func makeCBRERequest(url string) (*http.Response, error) {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
