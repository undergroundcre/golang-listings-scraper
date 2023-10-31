package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func omada() {
	url := "https://omada-cre.com/wp-json/wp/v2/listing?filter%5Borderby%5D=menu_order&filter%5Border%5D=ASC&per_page=100&filter[stype]=search&"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var listings []map[string]interface{}
	if err := json.Unmarshal(body, &listings); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	for _, listing := range listings {
		// Extract the parameters as before
		address := listing["_listing_address"]
		latitude := listing["_listing_latitude"]
		longitude := listing["_listing_longitude"]
		listingType := listing["_listing_listing_type"]
		var listingSize string
		if size, ok := listing["_listing_size_omada"].(string); ok {
			// Find the content inside <p> tags
			startTag := "<p>"
			endTag := "</p>"
			startIdx := strings.Index(size, startTag)
			endIdx := strings.Index(size, endTag)
			if startIdx != -1 && endIdx != -1 && endIdx > startIdx+len(startTag) {
				listingSize = size[startIdx+len(startTag) : endIdx]
			}
		}
		// Extract listingSize as before

		link := listing["link"]
		photo := listing["featured_image_src"]
		// Extract list price and sale price

		// Extract transaction type
		var transactionType string
		if ptiTaxonomies, ok := listing["pti_taxonomies"].(map[string]interface{}); ok {
			if status, ok := ptiTaxonomies["status"].([]interface{}); ok && len(status) > 0 {
				if statusItem, ok := status[0].(map[string]interface{}); ok {
					if transactionName, ok := statusItem["name"].(string); ok {
						transactionType = transactionName
					}
				}
			}
		}

		// Create a Scraper struct with the extracted data
		data := Scraper{
			URL:         fmt.Sprintf("%v", link),
			Asset:       fmt.Sprintf("%v", listingType),
			Transaction: transactionType,
			Location:    fmt.Sprintf("%v", address),
			Size:        listingSize,
			Latitude:    fmt.Sprintf("%v", latitude),
			Longitude:   fmt.Sprintf("%v", longitude),
			Photo:       fmt.Sprintf("%v", photo),
		}

		// Send the data to the datastore
		sendDataToDatastoreomada(data)
	}

	fmt.Println("Data from Omada sent to datastore successfully")
}

func sendDataToDatastoreomada(data Scraper) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}

	resp, err := http.Post("https://jsonserver-production-0d88.up.railway.app/add", "application/json", strings.NewReader(string(jsonData)))
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
