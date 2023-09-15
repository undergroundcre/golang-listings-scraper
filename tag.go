package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)


func extractPageURLsAndImages(body string) []Scraper {
	re := regexp.MustCompile(`image:\s+"([^"]+)",\s+page_url:\s+"([^"]+)"[\s\S]+?propertyType:\s+"([^"]*)?"[\s\S]+?price:\s+"([^"]*)?"`)
	matches := re.FindAllStringSubmatch(body, -1)

	var data []Scraper
	for _, match := range matches {
		dataItem := Scraper{
			URL:   "https://www.tag.ca" + match[2],
			Photo: "https://www.tag.ca" + match[1],
			Asset: match[3], // Assign propertyType to Asset
			Price: match[4], // Assign price to Price
		}

		data = append(data, dataItem)
	}

	return data
}

func tag() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.tag.ca/properties/", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.0.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Referer", "https://www.google.com/")
	req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"112\", \"Not_A Brand\";v=\"24\", \"Opera\";v=\"98\"")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"Linux\"")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		var body []byte
		for {
			chunk := make([]byte, 1024)
			n, err := reader.Read(chunk)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if n == 0 {
				break
			}
			body = append(body, chunk[:n]...)
		}

		data := extractPageURLsAndImages(string(body))

		for _, entry := range data {
			sendDataToServerz(entry)
		}

		log.Println("Scraping and sending data completed.")
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	}
}

func sendDataToServerz(data Scraper) {
	jsonData := fmt.Sprintf(`{"URL": "%s", "Transaction": "%s", "Photo": "%s", "Price": "%s"}`, data.URL, data.Asset, data.Photo, data.Price)
	resp, err := http.Post("https://jsonserver-production-799f.up.railway.app/add", "application/json", strings.NewReader(jsonData))
	if err != nil {
		log.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Data sent to the server")
}
