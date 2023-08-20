package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func cushman() {
	url := "https://cwedm.com/wp-admin/admin-ajax.php"
	headers := getHeaders()
	payload := getPayload()

	client := &http.Client{}
	payloadData := mapToURLValuez(payload)
	req, err := http.NewRequest("POST", url, strings.NewReader(payloadData.Encode()))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	var bodyReader io.Reader = resp.Body

	// Check if the response is gzipped and decompress if needed
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal("Error creating gzip reader:", err)
		}
		defer gzipReader.Close()

		bodyReader = gzipReader
	}

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		log.Fatal("Error decoding JSON response:", err)
	}

	results, ok := responseJSON["results"].([]interface{})
	if !ok {
		log.Fatal("Error parsing results from JSON response")
	}
	// Open the file in append mode
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	for _, result := range results {
		resultMap, isMap := result.(map[string]interface{})
		if !isMap {
			log.Println("Skipping invalid result entry")
			continue
		}

		writeResultToFile(file, resultMap)
	}
}

func getHeaders() map[string]string {
	return map[string]string{
		"authority":       "cwedm.com",
		"method":          "POST",
		"path":            "/wp-admin/admin-ajax.php",
		"scheme":          "https",
		"accept":          "application/json, text/javascript, */*; q=0.01",
		"accept-encoding": "gzip, deflate, br",
		"content-type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"origin":          "https://cwedm.com",
		"referer":         "https://cwedm.com/all-listings/",
		"sec-ch-ua":       "\"Chromium\";v=\"112\", \"Not_A Brand\";v=\"24\", \"Opera\";v=\"98\"",
		"sec-ch-ua-mobile": "?0",
		"sec-ch-ua-platform": "\"Linux\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.0.0",
		"x-requested-with":   "XMLHttpRequest",
	}
}

func getPayload() map[string]string {
	return map[string]string{
		"action":                       "rem_map_filters",
		"icons_by_meta":                "",
		"icons_data":                   "{\"\":{\"static\":\"\",\"hover\":\"\"}}",
		"radius":                       "",
		"radius_unit":                  "mi",
		"latitude":                     "",
		"longitude":                    "",
		"search_property":              "",
		"property_land_acres":          "",
		"range[property_building_sf][min]": "0.00",
		"range[property_building_sf][max]": "500,000.00",
		"range[property_building_sf][default_min]": "0",
		"range[property_building_sf][default_max]": "500000",
		"property_state_iso":           "",
		"price_min":                    "0.00",
		"price_max":                    "25,000,000.00",
		"price_min_default":            "0.00",
		"price_max_default":            "25,000,000.00",
		"range[lease_rate_sf][min]":    "0.00",
		"range[lease_rate_sf][max]":    "135.00",
		"range[lease_rate_sf][default_min]": "0",
		"range[lease_rate_sf][default_max]": "135.43",
		"range[lease_rate_per_acre][min]": "0.00",
		"range[lease_rate_per_acre][max]": "200.00",
		"range[lease_rate_per_acre][default_min]": "0",
		"range[lease_rate_per_acre][default_max]": "200",
		"agent_id":                     "",
		"order":                        "DESC",
		"orderby":                      "date",
	}
}

func mapToURLValuez(m map[string]string) url.Values {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, v)
	}
	return values
}

func writeResultToFile(file *os.File, result map[string]interface{}) {

	title := result["title"]
	address := result["address"].(string)
	latitude := result["latitude"]
	longitude := result["longitude"]
	propertyBoxHTML := result["property_box"].(string)

	// Parse the property_box HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(propertyBoxHTML))
	if err != nil {
		log.Println("Error parsing property_box HTML:", err)
		return
	}

	// Extract specific information from the parsed HTML
	area := strings.TrimSpace(doc.Find(".detail.inline-property-icons span[title='Area']").Text())

	// Extract image URL
	imageURL, _ := doc.Find(".img-container").Attr("style")
	imageURL = strings.TrimPrefix(strings.TrimSuffix(imageURL, "')"), "background-image:url('")

    // Extract information from the <a href> link
    link, _ := doc.Find(".rem-box-maps a").Attr("href")

   
	data := Scraper{
		URL:         link,
		Asset:       title.(string),
		Transaction: "",
		Location:    address,
		Size:        area,
		Latitude:    latitude.(string),
		Longitude:   longitude.(string),
		Photo:       imageURL,
	}

	// Send the data to the datastore
	sendDataToDatastoreloopnet(data)


}