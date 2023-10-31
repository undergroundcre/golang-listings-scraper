package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func cbre() {
	url := "https://www.cbre.ca/property-api/propertylistings/query?Site=ca-comm&RadiusType=Kilometers&CurrencyCode=CAD&Unit=sqft&Interval=Annually&Common.HomeSite=ca-comm&lon=-106.346771&Lat=56.130366&Lon=-122.583046&PolygonFilters=%5B%5B%2270%2C-50%22%2C%2242%2C-50%22%2C%2242%2C-142%22%2C%2270%2C-142%22%5D%5D&Common.Aspects=isLetting,isSale&PageSize=5000&Page=1&Sort=desc(Common.LastUpdated)&Common.UsageType=Office%2CRetail%2CIndustrial%2CLand%2CMultifamily%2CResidential%2CReligious%2COutdoorRecreational%2CIndoorRecreational%2CHotels%2CHealthcare%2CEducation%26Culture%2COpenStorage&Common.IsParent=true&_select=Dynamic.PrimaryImage,Common.ActualAddress,Common.Charges,Common.NumberOfBedrooms,Common.PrimaryKey,Common.UsageType,Common.Coordinate,Common.Aspects,Common.ListingCount,Common.IsParent,Common.HomeSite,Common.Agents,Common.PropertySubType,Common.ContactGroup,Common.Highlights,Common.Walkthrough,Common.MinimumSize,Common.MaximumSize,Common.TotalSize,Common.GeoLocation,Common.Sizes"
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

		// Create an array to hold the data to be sent to the datastore
		var dataToSend []Scraper

		for _, documentList := range documents {
			for _, document := range documentList.([]interface{}) {
				documentData := document.(map[string]interface{})

				primaryKey := documentData["Common.PrimaryKey"].(string)
				usageType := documentData["Common.UsageType"].(string)

				var transactionType string
				if aspects, ok := documentData["Common.Aspects"].([]interface{}); ok && len(aspects) > 0 {
					transactionType = aspects[0].(string)
				}
				var leaseprice string
				var salePrice string // Use an appropriate data type for the amount, e.g., int

				if charges, ok := documentData["Common.Charges"].([]interface{}); ok {
					if len(charges) > 2 { // Ensure that there's an element with index 2
						chargeData, isMap := charges[2].(map[string]interface{})
						if isMap {
							if chargeKind, exists := chargeData["Common.ChargeKind"].(string); exists {
								if chargeKind == "SalePrice" {
									if amountData, hasAmount := chargeData["Common.Amount"].(float64); hasAmount {
										salePrice = fmt.Sprintf("%.2f", amountData) // Convert to string with 2 decimal places
									}
								} else {
									if chargeAmount, exists := chargeData["Common.Amount"].(float64); exists {
										leaseprice = fmt.Sprintf("%.2f", chargeAmount) // Convert to string with 2 decimal places
									}
								}
							}
						}
					}
				}

				var propertytype string
				if transactionType == "isSale" {
					propertytype = "Sale"
				} else if transactionType == "isLetting" {
					propertytype = "Lease"
				}

				var state string
				if statez, ok := documentData["Common.ActualAddress"].(map[string]interface{}); ok {
					state, _ = statez["Common.Region"].(string)
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

				constructedURL := fmt.Sprintf("https://www.cbre.ca/properties/commercial-space/details/%s", primaryKey)

				// Add the scraped data to the array
				dataToSend = append(dataToSend, Scraper{
					URL:         constructedURL,
					Asset:       usageType,
					Transaction: propertytype,
					Location:    line1,
					Size:        fmt.Sprintf("%.2f %s", size, units),
					Latitude:    fmt.Sprintf("%f", lat),
					Longitude:   fmt.Sprintf("%f", lon),
					Photo:       sourceUri,
					State:       state,
					LeaseRate:   leaseprice,
					Price:       salePrice,
				})
			}
		}

		// Send the scraped data to the datastore
		for _, data := range dataToSend {
			sendDataToDatastorecbre(data)
		}

		fmt.Println("Data sent to datastore successfully")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
}

func sendDataToDatastorecbre(data Scraper) {
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
