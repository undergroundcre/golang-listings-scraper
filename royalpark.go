package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
	"golang.org/x/net/html"
)

func royalpark() {
	url := "https://royalparkrealty.com/wp-admin/admin-ajax.php?a=40009.587902062645"

	payload := "action=getProperties&search[team_member]=&page=1&max=8&view=Map%20View"
	headers := map[string]string{
		"authority":            "royalparkrealty.com",
		"method":               "POST",
		"path":                 "/wp-admin/admin-ajax.php?a=40009.587902062645",
		"scheme":               "https",
		"accept":               "application/json, text/javascript, */*; q=0.01",
		"accept-encoding":      "gzip, deflate, br",
		"content-length":       "72",
		"content-type":         "application/x-www-form-urlencoded; charset=UTF-8",
		"cookie":               "saved_search_params=%7B%22name_id%22%3A%22%22%2C%22transaction_type%22%3A%22%22%2C%22property_type%22%3A%22%22%2C%22location%22%3A%22%22%2C%22building_size%22%3A%22%22%2C%22site_size%22%3A%22%22%2C%22net_lease_price%22%3A%22%22%2C%22sale_price%22%3A%22%22%2C%22property_sort_by%22%3A%22%22%2C%22property_sort_type%22%3A%22ASC%22%2C%22team_member%22%3A%22%22%7D; _gid=GA1.2.714389502.1691620324; __utmc=26616382; _ga=GA1.2.1441758266.1691239157; _ga_ZWP1JJSGW7=GS1.1.1691705738.8.0.1691705738.0.0.0; _ga_6V0WB7RCDX=GS1.1.1691705738.8.0.1691705738.0.0.0; _gat=1; __utma=26616382.1441758266.1691239157.1691700604.1691705748.7; __utmz=26616382.1691705748.7.5.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmt=1; __utmb=26616382.1.10.1691705748",
		"origin":               "https://royalparkrealty.com",
		"referer":              "https://royalparkrealty.com/property-search/",
		"sec-ch-ua":            "\"Chromium\";v=\"112\", \"Not_A Brand\";v=\"24\", \"Opera\";v=\"98\"",
		"sec-ch-ua-mobile":     "?0",
		"sec-ch-ua-platform":   "\"Linux\"",
		"sec-fetch-dest":       "empty",
		"sec-fetch-mode":       "cors",
		"sec-fetch-site":       "same-origin",
		"user-agent":           "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.0.0",
		"x-requested-with":     "XMLHttpRequest",
	}
	

	response, err := sendPostRequest(url, payload, headers)
	if err != nil {
		log.Fatal("Error sending POST request:", err)
	}

	// Extract relevant information from JSON response
	properties, err := extractProperties(response)
	if err != nil {
		log.Fatal("Error extracting properties:", err)
	}

	// Write information to the file
	err = writePropertiesToFile(properties)
	if err != nil {
		log.Fatal("Error writing properties to file:", err)
	}
}


func sendPostRequest(url, payload string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response is compressed with Brotli
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "br":
		reader = brotli.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return body, nil
}


func extractProperties(response []byte) ([]map[string]interface{}, error) {
	// Parse the JSON response
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(response, &jsonResponse); err != nil {
		return nil, err
	}

	// Extract properties from the response
	properties, ok := jsonResponse["properties"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("properties not found in JSON response")
	}

	var extractedProperties []map[string]interface{}
	for _, prop := range properties {
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			continue
		}
		extractedProperties = append(extractedProperties, propMap)
	}

	return extractedProperties, nil
}

func writePropertiesToFile(properties []map[string]interface{}) error {
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()
// ...

for _, prop := range properties {
	address, _ := prop["address"].(string)
	size := html.UnescapeString(prop["size"].(string))
    price := html.UnescapeString(prop["price"].(string))
	latitude, _ := prop["latitude"].(float64)
	longitude, _ := prop["longitude"].(float64)
	permalink, _ := prop["permalink"].(string)
	transactionType, _ := prop["transaction_type"].(string)
	thumbnail, _ := prop["thumbnail"].(string)
	thumbnailSrc := extractThumbnailSrc(thumbnail)


	// Write property information to file
	file.WriteString("Address: " + address + "\n")
	file.WriteString("Size: " + size + "\n")
	file.WriteString("Price: " + price + "\n")
	file.WriteString(fmt.Sprintf("Latitude: %f\n", latitude))
	file.WriteString(fmt.Sprintf("Longitude: %f\n", longitude))
	file.WriteString("Permalink: " + permalink + "\n")
	file.WriteString("Transaction Type: " + transactionType + "\n")
	file.WriteString("Thumbnail: " + thumbnailSrc + "\n")
	file.WriteString("-------------------------------\n")
}

// ...


	return nil
}

func extractThumbnailSrc(thumbnailContent string) string {
    doc, err := html.Parse(strings.NewReader(thumbnailContent))
    if err != nil {
        return ""
    }
    var extractSrc func(*html.Node) string
    extractSrc = func(n *html.Node) string {
        if n.Type == html.ElementNode && n.Data == "img" {
            for _, attr := range n.Attr {
                if attr.Key == "src" {
                    return attr.Val
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            if src := extractSrc(c); src != "" {
                return src
            }
        }
        return ""
    }
    return extractSrc(doc)
}