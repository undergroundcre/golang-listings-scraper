package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

		for _, placard := range placardElements.Nodes {
			headerLink := goquery.NewDocumentFromNode(placard).Find(".placard-content a")
			headerURL, _ := headerLink.Attr("href")

			nameFull := goquery.NewDocumentFromNode(placard).Find(".header-col h4")
			name := nameFull.Text()

			nameFulltwo := goquery.NewDocumentFromNode(placard).Find(".header-col h6")
			nametwo := nameFulltwo.Text()

			locationFull := goquery.NewDocumentFromNode(placard).Find(".header-col a")
			location := locationFull.Text()

			priceElement := goquery.NewDocumentFromNode(placard).Find(".placard-info .data-points-2c li[name='Price']")
			price := strings.TrimSpace(priceElement.Text())

			assetTypeElement := goquery.NewDocumentFromNode(placard).Find(".placard-info .data-points-2c li:nth-child(2)")
			assetType := strings.TrimSpace(assetTypeElement.Text())

			slide := goquery.NewDocumentFromNode(placard).Find(".slide.active")
			imgSrc, _ := slide.Find("img").Attr("src")

			result := fmt.Sprintf("URL: %s\nName: %s\nPrice: %s\nLocation: %s\nAsset type: %s\nTransaction Type: Lease\nImage: %s\n-------------------------------\n",
				headerURL, name + " " + nametwo, price, location ,assetType, imgSrc)

			result = strings.TrimSpace(result) + "\n" // Trim spaces and add a newline
			_, err := file.WriteString(result)
			if err != nil {
				log.Fatal(err)
			}
		}

		page++
	}

}
