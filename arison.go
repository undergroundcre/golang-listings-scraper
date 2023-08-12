package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func fetchPageData(url string, payload map[string]string, headers map[string]string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(mapToURLValues(payload).Encode()))
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

func arison() {
	url := "https://buildout.com/plugins/3e7f33c297c67cca1e29a9cb4e99b1ed75f8efd0/inventory"

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
		"utf8":                                    "✓",
		"polygon_geojson":                         "",
		"lat_min":                                 "",
		"lat_max":                                 "",
		"lng_min":                                 "",
		"lng_max":                                 "",
		"mobile_lat_min":                          "",
		"mobile_lat_max":                          "",
		"mobile_lng_min":                          "",
		"mobile_lng_max":                          "",
		"page":                                    "0",
		"map_display_limit":                       "5000",
		"map_type":                                "roadmap",
		"country_restrictions":                    "ca",
		"custom_map_marker_url":                   "/s3.amazonaws.com/buildout-production/brandings/6929/profile_photo/small.png?1600277316",
		"use_marker_clusterer":                    "true",
		"placesAutoComplete":                      "",
		"q[type_use_offset_eq_any][]":             "",
		"q[sale_or_lease_eq]":                     "",
		"q[state_eq_any][]":                       "",
		"q[listings_data_max_space_available_on_market_gteq]": "",
		"q[listings_data_min_space_available_on_market_lteq]": "",
		"q[max_lease_rate_gteq]":                  "",
		"q[min_lease_rate_lteq]":                  "",
		"q[max_lease_rate_monthly_gteq]":          "",
		"q[min_lease_rate_monthly_lteq]":          "",
		"q[has_broker_ids][]":                    "",
		"q[s][]":                                 "max_lease_rate desc",
	}

	file, err := os.Create("data.txt")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	pageNumber := 0
	for {
		payload["page"] = fmt.Sprintf("%d", pageNumber)
		data, err := fetchPageData(url, payload, headers)
		if err != nil {
			log.Printf("Failed to fetch data for page %d: %s", pageNumber, err)
			break
		}

		for _, entry := range data["inventory"].([]interface{}) {
			entryMap := entry.(map[string]interface{})
			url := entryMap["also_for_sale_or_lease_url"].(string)
			if url != "" {
				url = strings.Replace(url, "sale", "lease", -1) // Replace all occurrences of "sale" with "lease"
				assetType := entryMap["property_sub_type_name"].(string)
				saleorlease := entryMap["sale"]
				location := entryMap["address"].(string)
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
		
				// Extract and prepare index_attributes
				attributes := ""
				if attrs, ok := entryMap["index_attributes"].([]interface{}); ok {
					for _, attr := range attrs {
						if attrSlice, isSlice := attr.([]interface{}); isSlice && len(attrSlice) == 2 {
							key := attrSlice[0].(string)
							value := attrSlice[1].(string)
							attributes += fmt.Sprintf("%s: %s\n", key, value)
						}
					}
				}
		
				file.WriteString(fmt.Sprintf("URL: %s\nAsset Type: %s\nTransaction Type: %s\nLocation: %s\nSize: %s\nLatitude: %s\nLongitude: %s\nPhoto: %s\n%s--------------------\n", url, assetType, transactiontype, location, size, LatitudeStr, LongitudeStr, photo, attributes))
			}
		}
		

		pageNumber++

		// Check if there are more pages to fetch
		if len(data["inventory"].([]interface{})) == 0 {
			break
		}
	}

}

func mapToURLValues(m map[string]string) url.Values {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, v)
	}
	return values
}
