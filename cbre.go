package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func cbre() {
	url := "https://www.cbre.ca/property-api/propertylistings/query?Site=ca-comm&RadiusType=Kilometers&CurrencyCode=CAD&Unit=sqft&Interval=Annually&Common.HomeSite=ca-comm&lon=-106.346771&Lat=56.130366&Lon=-122.583046&PolygonFilters=%5B%5B%2270%2C-50%22%2C%2242%2C-50%22%2C%2242%2C-142%22%2C%2270%2C-142%22%5D%5D&Common.Aspects=isSale&PageSize=400&Page=1&Sort=desc(Common.LastUpdated)&Common.UsageType=Office&Common.IsParent=true&_select=Dynamic.PrimaryImage,Common.ActualAddress,Common.Charges,Common.NumberOfBedrooms,Common.PrimaryKey,Common.UsageType,Common.Coordinate,Common.Aspects,Common.ListingCount,Common.IsParent,Common.HomeSite,Common.Agents,Common.PropertySubType,Common.ContactGroup,Common.Highlights,Common.Walkthrough,Common.MinimumSize,Common.MaximumSize,Common.TotalSize,Common.GeoLocation,Common.Sizes"
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var jsonResponse map[string]interface{}
		if err := json.Unmarshal(body, &jsonResponse); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		documents := jsonResponse["Documents"].([]interface{})

		// Open the file in append mode
		dataFile, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer dataFile.Close()

		for _, documentList := range documents {
			for _, document := range documentList.([]interface{}) {
				documentData := document.(map[string]interface{})

				primaryKey := documentData["Common.PrimaryKey"].(string)
				usageType := documentData["Common.UsageType"].(string)

				var transactionType string
				if aspects, ok := documentData["Common.Aspects"].([]interface{}); ok && len(aspects) > 0 {
					transactionType = aspects[0].(string)
				}

				var lon, lat float64
				if coordinate, ok := documentData["Common.Coordinate"].(map[string]interface{}); ok {
					lon, _ = coordinate["lon"].(float64)
					lat, _ = coordinate["lat"].(float64)
				}

				var size float64
				var units string
				if totalSizeArr, ok := documentData["Dynamic.TotalSize"].([]interface{}); ok && len(totalSizeArr) > 0 {
					totalSize := totalSizeArr[0].(map[string]interface{})
					size, _ = totalSize["Common.Size"].(float64)
					units, _ = totalSize["Common.Units"].(string)
				}

				var sourceUri, line1 string
				if primaryImage, ok := documentData["Dynamic.PrimaryImage"].(map[string]interface{}); ok {
					if imageResource, ok := primaryImage["Common.ImageResources"].([]interface{}); ok && len(imageResource) > 0 {
						sourceUri = imageResource[0].(map[string]interface{})["Source.Uri"].(string)
					}
				}
				if actualAddressData, ok := documentData["Common.ActualAddress"].(map[string]interface{}); ok {
					if postalAddresses, ok := actualAddressData["Common.PostalAddresses"].([]interface{}); ok && len(postalAddresses) > 1 {
						secondAddress := postalAddresses[1].(map[string]interface{})
						if line1Value, line1Exists := secondAddress["Common.Line1"].(string); line1Exists {
							line1 = line1Value
						}
					}
				}

				constructedURL := fmt.Sprintf("https://www.cbre.ca/property-details/%s", primaryKey)

				// Write extracted data and constructed URL to the file
				_, err := dataFile.WriteString(fmt.Sprintf("Name: %s\nProperty Use: %s\nTransaction Type: %s\nLongitude: %f\nLatitude: %f\nSize: %.2f %s\nSecond Address Line1: %s\nPhoto: %s\nConstructed URL: %s\n-------------------------------------\n",
					primaryKey, usageType, transactionType, lon, lat, size, units, line1, sourceUri, constructedURL))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		fmt.Println("Parameters extracted and appended to data.txt")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
}
