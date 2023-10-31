package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func loopnet() {
	baseURL := "https://www.loopnet.ca/search/commercial-real-estate/canada/for-lease/"
	page := 1
	hasListings := true

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for hasListings {
		url := fmt.Sprintf("%s%d/", baseURL, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")

		response, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		placardElements := doc.Find(".placard")
		if placardElements.Length() == 0 {
			hasListings = false
			break
		}
		canadianStates := []string{"AB", "BC", "MB", "NB", "NL", "NS", "NT", "NU", "ON", "PE", "QC", "SK", "YT"}

		for _, placard := range placardElements.Nodes {
			priceElement := goquery.NewDocumentFromNode(placard).Find("li[name=Price]")
			var price string
			if priceElement.Length() > 0 {
				price = strings.TrimSpace(priceElement.Text())
			}
			headerLink := goquery.NewDocumentFromNode(placard).Find(".placard-content a")
			headerURL, _ := headerLink.Attr("href")

			locationFull := goquery.NewDocumentFromNode(placard).Find(".header-col a")
			location := locationFull.Text()

			state := ""
			for _, canadianState := range canadianStates {
				if strings.Contains(location, canadianState) {
					state = canadianState
					break
				}
			}

			assetTypeElement := goquery.NewDocumentFromNode(placard).Find(".placard-info .data-points-2c li:nth-child(2)")
			assetType := strings.TrimSpace(assetTypeElement.Text())
			assetTypeParts := strings.SplitN(assetType, "SF", 2)

			var size, usagetype string
			if len(assetTypeParts) > 1 {
				size = strings.TrimSpace(assetTypeParts[0] + "SF")
				usagetype = strings.TrimSpace(assetTypeParts[1])
			} else {
				size = ""
				usagetype = strings.TrimSpace(assetType)
			}

			slide := goquery.NewDocumentFromNode(placard).Find(".slide.active")
			imgSrc, _ := slide.Find("img").Attr("src")

			constructedURL := headerURL
			transactiontype := "Lease"
			line1 := location
			lat := "" // No latitude data in this source
			lon := "" // No longitude data in this source
			sourceuri := imgSrc

			// Create a Scraper struct with the extracted data
			data := Scraper{
				URL:         constructedURL,
				Asset:       usagetype,
				Transaction: transactiontype,
				Location:    line1,
				Size:        size,
				Latitude:    lat,
				Longitude:   lon,
				Photo:       sourceuri,
				LeaseRate:   price,
				State:       state,
			}

			// Send the data to the datastore
			sendDataToDatastoreloopnet(data)
		}

		page++
	}

	fmt.Println("Data from LoopNet sent to datastore successfully")
}

func sendDataToDatastoreloopnet(data Scraper) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}

	resp, err := http.Post("https://jsonserver-production-799f.up.railway.app/add", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("Failed to send data to datastore:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Data sent to datastore:", string(body))
}
