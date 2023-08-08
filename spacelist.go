package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ListingData struct {
	URL        string
	ImageSrc   string
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

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	failedFile, err := os.Create("failed_scrape.txt") // Create a file for failed URLs
	if err != nil {
		log.Fatal("Error creating failed scrape file:", err)
	}
	defer failedFile.Close()

	// Add a ticker to track the time and determine when to pause the fetching
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case listing, ok := <-listingsChan:
			if !ok {
				listingsChan = nil // Set to nil to avoid sending data to a closed channel
			} else if listing.Name == "" || len(listing.KeyBoldMap) == 0 {
				fmt.Fprintln(failedFile, listing.URL) // Write failed URL to the failed_scrape.txt file
			} else {
				fmt.Fprintf(file, "URL: %s\n", listing.URL)
				fmt.Fprintf(file, "Location: %s\n", listing.Name)
				fmt.Fprintf(file, "Photo: %s\n", listing.ImageSrc)
				for key, boldFont := range listing.KeyBoldMap {
					fmt.Fprintf(file, "%s | %s\n", key, boldFont)
				}
				fmt.Fprintln(file, strings.Repeat("-", 20))
			}

		case failedURL, ok := <-failedURLsChan:
			if !ok {
				failedURLsChan = nil // Set to nil to avoid sending data to a closed channel
			} else {
				fmt.Fprintln(failedFile, failedURL) // Write failed URL to the failed_scrape.txt file
			}

		case <-ticker.C: // After every 10 minutes, pause and wait for a random duration
			pauseDuration := 3 + rand.Intn(4) // Generates a random value between 3 to 6
			pauseTime := time.Now()
			fmt.Printf("[%s] Pausing fetching for %d minutes...\n", pauseTime.Format(time.RFC3339), pauseDuration)
			time.Sleep(time.Duration(pauseDuration) * time.Minute)
			resumeTime := time.Now()
			fmt.Printf("[%s] Resuming fetching...\n", resumeTime.Format(time.RFC3339))

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
			failedURLsChan <- currentURL // Send failed URL to the channel
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println("Error parsing HTML:", err)
			failedURLsChan <- currentURL // Send failed URL to the channel
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
				failedURLsChan <- fullURL // Send failed URL to the channel
				continue
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Println("Error parsing HTML:", err)
				failedURLsChan <- fullURL // Send failed URL to the channel
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

			// Scrape image src
			var imageSrc string
			doc.Find("a.listing-image").Each(func(_ int, s *goquery.Selection) {
				if imageSrc == "" { // Only get the first image
					imageSrc, _ = s.Attr("href")
				}
			})

			listingsChan <- ListingData{
				URL:        fullURL,
				ImageSrc:   imageSrc, // Add the image source to the struct
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

