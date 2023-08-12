package main

import (
	"fmt"
	"log"
	"net/http"
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

		headerElements := doc.Find(".header-col.header-left")
		sizeLocElements := doc.Find(".header-col.header-right.text-right")
		slideElements := doc.Find(".slide.active")

		if headerElements.Length() == 0 || sizeLocElements.Length() == 0 || slideElements.Length() == 0 {
			hasListings = false
			break
		}

		headerElements.Each(func(index int, div *goquery.Selection) {
			leftH4Text := div.Find("a.left-h4").Text()
			leftH6Text := div.Find("a.left-h6").Text()

			headerLink := div.Find("a.left-h4")
			headerURL, _ := headerLink.Attr("href")

			locationElement := sizeLocElements.Eq(index).Find("a.right-h6")
			location := locationElement.Text()

			rightH4Element := sizeLocElements.Eq(index).Find("a.right-h4")
			rightH4Text := rightH4Element.Text()

			slide := slideElements.Eq(index)
			imgSrc, _ := slide.Find("img").Attr("src")

			fmt.Println("URL:", headerURL)
			fmt.Println("Name:", leftH4Text, leftH6Text)
			fmt.Println("Type:", rightH4Text)
			fmt.Println("Location:", location)
			fmt.Println("Image:", imgSrc)
			fmt.Println("-------------------------------")
		})

		page++
	}

	fmt.Println("No more listings found.")
}
