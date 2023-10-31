package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func fetchPageDataYeg(url string, payload map[string]string, headers map[string]string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(mapToURLValuesYeg(payload)))
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func yeg() {
	url := "https://buildout.com/plugins/df37c8376797e16bea8783f81133fd37e2450e3f/inventory"

	headers := map[string]string{
		"accept":             "application/json, text/javascript, */*; q=0.01",
		"content-type":       "application/x-www-form-urlencoded; charset=UTF-8",
		"sec-ch-ua":          "\"Chromium\";v=\"112\", \"Not_A Brand\";v=\"24\", \"Opera\";v=\"98\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Linux\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"x-newrelic-id":      "Vg4GU1RRGwIJUVJUAwY=",
		"x-requested-with":   "XMLHttpRequest",
	}

	payload := map[string]string{
		"utf8":                        "✓",
		"polygon_geojson":             "",
		"lat_min":                     "",
		"lat_max":                     "",
		"lng_min":                     "",
		"lng_max":                     "",
		"mobile_lat_min":              "",
		"mobile_lat_max":              "",
		"mobile_lng_min":              "",
		"mobile_lng_max":              "",
		"page":                        "0",
		"map_display_limit":           "5000",
		"map_type":                    "roadmap",
		"country_restrictions":        "ca",
		"custom_map_marker_url":       "/s3.amazonaws.com/buildout-production/brandings/6929/profile_photo/small.png?1600277316",
		"use_marker_clusterer":        "true",
		"placesAutoComplete":          "",
		"q[type_use_offset_eq_any][]": "",
		"q[sale_or_lease_eq]":         "",
		"q[state_eq_any][]":           "",
		"q[listings_data_max_space_available_on_market_gteq]": "",
		"q[listings_data_min_space_available_on_market_lteq]": "",
		"q[max_lease_rate_gteq]":                              "",
		"q[min_lease_rate_lteq]":                              "",
		"q[max_lease_rate_monthly_gteq]":                      "",
		"q[min_lease_rate_monthly_lteq]":                      "",
		"q[has_broker_ids][]":                                 "",
		"q[s][]":                                              "sale_price desc",
	}

	pageNumber := 0
	for {
		payload["page"] = fmt.Sprintf("%d", pageNumber)
		data, err := fetchPageDataYeg(url, payload, headers)
		if err != nil {
			log.Printf("Failed to fetch data for page %d: %s", pageNumber, err)
			break
		}

		inventory, inventoryExists := data["inventory"].([]interface{})
		if !inventoryExists || len(inventory) == 0 {
			log.Printf("No more listings. Exiting...")
			break // Break out of the loop when there are no more listings
		}

		for _, entry := range inventory {
			entryMap := entry.(map[string]interface{})
			url := entryMap["show_link"].(string)
			if url != "" {
				assetType := entryMap["property_sub_type_name"].(string)
				saleorlease := entryMap["sale"]
				state := entryMap["state"].(string)
				location := entryMap["address"].(string)
				indexAttributes, ok := entryMap["index_attributes"].([]interface{})
				if !ok {
					fmt.Println(ok)
				}
				var leaseRate string
				// Loop through the index_attributes to find the "Lease Rate" value
				for _, attr := range indexAttributes {
					if attrArray, isArray := attr.([]interface{}); isArray && len(attrArray) == 2 {
						if label, isString := attrArray[0].(string); isString && label == "Lease Rate" {
							if rate, isString := attrArray[1].(string); isString {
								leaseRate = rate
								break
							}
						}
					}
				}

				var price string

				// Access the index_attributes array
				indexAttributez, okz := entryMap["index_attributes"].([]interface{})
				if !okz {
					fmt.Println(okz)
				}

				// Loop through the index_attributes to find the "Price" value
				for _, attr := range indexAttributez {
					if attrArray, isArray := attr.([]interface{}); isArray && len(attrArray) == 2 {
						if label, isString := attrArray[0].(string); isString && label == "Price" {
							if value, isString := attrArray[1].(string); isString {
								price = value
								break
							}
						}
					}
				}

				size := entryMap["size_summary"].(string)
				photo := entryMap["photo_url"].(string)
				Latitude, latOk := entryMap["latitude"].(float64)
				Longitude, longOk := entryMap["longitude"].(float64)

				var LatitudeStr, LongitudeStr string
				if latOk && longOk {
					LatitudeStr = strconv.FormatFloat(Latitude, 'f', 6, 64)
					LongitudeStr = strconv.FormatFloat(Longitude, 'f', 6, 64)
				} else {
					LatitudeStr = "N/A"
					LongitudeStr = "N/A"
				}

				var transactiontype string
				if saleorlease == true {
					transactiontype = "Sale"
				} else {
					transactiontype = "Lease"
				}

				var leasert string
				if leaseRate != "" {
					leasert = leaseRate
				} else {
					leasert = ""
				}

				scraperData := Scraper{
					URL:         url,
					Asset:       assetType,
					Transaction: transactiontype,
					Location:    location,
					Size:        size,
					Latitude:    LatitudeStr,
					Longitude:   LongitudeStr,
					Photo:       photo,
					Price:       price,
					LeaseRate:   leasert,
					State:       state,
				}

				// Send data to datastore
				sendDataToDatastoreYeg(scraperData)
			}
		}

		pageNumber++
	}

	log.Printf("Script finished successfully.")
}

func mapToURLValuesYeg(m map[string]string) string {
	values := []string{}
	for k, v := range m {
		values = append(values, k+"="+url.QueryEscape(v))
	}
	return strings.Join(values, "&")
}

func sendDataToDatastoreYeg(data Scraper) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %s", err)
		return
	}

	resp, err := http.Post("https://jsonserver-production-0d88.up.railway.app/add", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		log.Printf("Failed to send data to datastore: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return
	}

	log.Printf("Data sent to datastore: %s", string(body))
}
