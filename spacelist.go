package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ListingData struct {
	URL        string
	Photo      string
	Name       string
	KeyBoldMap map[string]string
}

func ScrapeListingsFromMainURLs() {
	mainURLs := []string{
		"https://e85.spacelist.ca/listings",
		"https://www.spacelist.ca/listings/ab",
		"https://e149.spacelist.ca/listings",
	}

	var wg sync.WaitGroup
	listingsChan := make(chan ListingData)
	failedURLsChan := make(chan string) // Channel to hold failed URLs

	for _, mainURL := range mainURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			ScrapeListings(url, listingsChan, failedURLsChan)
		}(mainURL)
	}

	go func() {
		wg.Wait()
		close(listingsChan)
		close(failedURLsChan)
	}()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case listing, ok := <-listingsChan:
			if !ok {
				listingsChan = nil // Set to nil to avoid sending data to a closed channel
			} else if listing.Name == "" || len(listing.KeyBoldMap) == 0 {
				log.Println("Failed to scrape data:", listing.URL)
			} else {
				sendDataToServer(listing) // Send the listing data to the server
			}

		case failedURL, ok := <-failedURLsChan:
			if !ok {
				failedURLsChan = nil // Set to nil to avoid sending data to a closed channel
			} else {
				log.Println("Failed to scrape data:", failedURL)
			}

		case <-ticker.C: // Pause for 2 minutes every 10 minutes
			log.Println("Pausing for 7 minutes...")
			time.Sleep(7 * time.Minute)
			log.Println("Resuming fetching...")

			if listingsChan == nil && failedURLsChan == nil {
				return // Exit the loop if both channels are closed
			}
		}

		if listingsChan == nil && failedURLsChan == nil {
			break // Exit the loop if both channels are closed
		}
	}
}

func ScrapeListings(url string, listingsChan chan ListingData, failedURLsChan chan string) {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	page := 1
	for {
		currentURL := fmt.Sprintf("%s/page/%d", url, page)

		resp, err := makeRequest(currentURL, headers)
		if err != nil {
			log.Println("Error making request:", err)
			failedURLsChan <- currentURL
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println("Error parsing HTML:", err)
			failedURLsChan <- currentURL
			return
		}

		var hrefs []string
		doc.Find("a.listing-card").Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			hrefs = append(hrefs, href)
		})

		if len(hrefs) == 0 {
			break // No more pages
		}

		for _, href := range hrefs {
			fullURL := href

			resp, err := makeRequest(fullURL, headers)
			if err != nil {
				log.Println("Error making request:", err)
				failedURLsChan <- fullURL
				continue
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Println("Error parsing HTML:", err)
				failedURLsChan <- fullURL
				continue
			}

			var name string
			doc.Find("h1.large-font").Each(func(_ int, s *goquery.Selection) {
				name = s.Text()
			})

			keyBoldMap := make(map[string]string)

			var keyText string
			doc.Find("span.key, span.bold-font").Each(func(_ int, s *goquery.Selection) {
				if strings.TrimSpace(s.Text()) != "" && strings.TrimSpace(s.Text()) != "Download Brochures" {
					if s.HasClass("key") {
						keyText = s.Text()
					} else if s.HasClass("bold-font") {
						keyBoldMap[keyText] = s.Text()
					}
				}
			})

			var imageSrc string
			doc.Find("a.listing-image").Each(func(_ int, s *goquery.Selection) {
				if imageSrc == "" {
					imageSrc, _ = s.Attr("href")
				}
			})

			listingsChan <- ListingData{
				URL:        fullURL,
				Photo:      imageSrc,
				Name:       name,
				KeyBoldMap: keyBoldMap,
			}
		}

		page++
	}
}

func makeRequest(url string, headers map[string]string) (*http.Response, error) {
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

func sendDataToServer(listing ListingData) {
	jsonData, err := json.Marshal(listing)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	resp, err := http.Post("https://spacelistmiddlescript-production.up.railway.app/data", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		log.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Data sent to server")
}

