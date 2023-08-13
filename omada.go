package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func omada() {
	url := "https://omada-cre.com/wp-json/wp/v2/listing?filter%5Borderby%5D=menu_order&filter%5Border%5D=ASC&per_page=100&filter[stype]=search&"

	// Set up the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add headers
	req.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
	// Add other headers...

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Unmarshal the JSON response
	var listings []map[string]interface{}
	if err := json.Unmarshal(body, &listings); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Create or open the output file
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	// Iterate through the listings and extract the desired parameters
	// Iterate through the listings and extract the desired parameters
	for _, listing := range listings {
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
		link := listing["link"]
		photo := listing["featured_image_src"]
			// Extract List Price and its value using string manipulation
			if priceValue, ok := listing["_listing_feature_omada"].(string); ok {
				listPriceTag := "<strong>List Price:</strong>"
				salePriceTag := "<strong>Sale Price:</strong>"
	
				// Find the index of List Price and Sale Price tags
				listPriceIndex := strings.Index(priceValue, listPriceTag)
				salePriceIndex := strings.Index(priceValue, salePriceTag)
	
				// Extract and write List Price value if found
				if listPriceIndex != -1 {
					listPriceEndIndex := strings.Index(priceValue[listPriceIndex+len(listPriceTag):], "</p>")
					if listPriceEndIndex != -1 {
						listPrice := strings.TrimSpace(priceValue[listPriceIndex+len(listPriceTag) : listPriceIndex+len(listPriceTag)+listPriceEndIndex])
						file.WriteString(fmt.Sprintf("Price: %s\n", listPrice))
					}
				}
	
				// Extract and write Sale Price value if found
				if salePriceIndex != -1 {
					salePriceEndIndex := strings.Index(priceValue[salePriceIndex+len(salePriceTag):], "</p>")
					if salePriceEndIndex != -1 {
						salePrice := strings.TrimSpace(priceValue[salePriceIndex+len(salePriceTag) : salePriceIndex+len(salePriceTag)+salePriceEndIndex])
						file.WriteString(fmt.Sprintf("Price: %s\n", salePrice))
					}
				}
			}

		// Extract nested parameter value
		if ptiTaxonomies, ok := listing["pti_taxonomies"].(map[string]interface{}); ok {
			if status, ok := ptiTaxonomies["status"].([]interface{}); ok && len(status) > 0 {
				if statusItem, ok := status[0].(map[string]interface{}); ok {
					if transactionType, ok := statusItem["name"].(string); ok {
						// Write the extracted transaction type to the file
						file.WriteString(fmt.Sprintf("Transaction Type: %v\n", transactionType))
					}
				}
			}
		}

		// Write the remaining extracted parameters to the file
		file.WriteString(fmt.Sprintf("Address: %v\n", address))
		file.WriteString(fmt.Sprintf("Latitude: %v\n", latitude))
		file.WriteString(fmt.Sprintf("Longitude: %v\n", longitude))
		file.WriteString(fmt.Sprintf("Listing Type: %v\n", listingType))
		file.WriteString(fmt.Sprintf("Photo: %v\n", photo))
		file.WriteString(fmt.Sprintf("URL: %v\n", link))
		file.WriteString(fmt.Sprintf("Size: %v\n", listingSize))
		file.WriteString("-------------------------------\n")
	}

}

