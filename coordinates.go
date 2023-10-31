package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func locationToCoordinate() {
	dsn := "root:CBCFeeEae6b-52aCf15g4EBAcA4Bcd-6@tcp(monorail.proxy.rlwy.net:26586)/scraper?timeout=60s"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Scraper{})
	// Query the database to fetch rows with empty latitude and longitude
	var scrapers []Scraper
	db.Where("latitude = '' AND longitude = ''").Find(&scrapers)

	// Iterate through the scrapers
	for _, scraper := range scrapers {
		// Send a request to the Google Maps API to geocode the location
		geocodeURL := "https://maps.googleapis.com/maps/api/geocode/json"
		apiKey := "AIzaSyAYjB1P6H3zXAU-ETAExNvJ4iXE47lXtAI" // Replace with your API key
		location := scraper.Location
		fmt.Println(location)

		// Create a request with custom headers
		req, err := http.NewRequest("GET", geocodeURL, nil)
		if err != nil {
			log.Println("Error creating geocoding request:", err)
			continue
		}

		// Add headers to the request
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		// Set query parameters
		q := req.URL.Query()
		q.Add("address", location)
		q.Add("key", apiKey)
		req.URL.RawQuery = q.Encode()

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending geocoding request:", err)
			continue
		}
		defer resp.Body.Close()

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading geocoding response:", err)
			continue
		}

		// Parse the response data
		var geocodingResult map[string]interface{}
		err = json.Unmarshal(body, &geocodingResult)
		if err != nil {
			log.Println("Error parsing geocoding response:", err)
			continue
		}

		// Extract the coordinates from the response and update the database
		if status, ok := geocodingResult["status"].(string); ok && status == "OK" {
			if results, ok := geocodingResult["results"].([]interface{}); ok && len(results) > 0 {
				if geometry, ok := results[0].(map[string]interface{})["geometry"].(map[string]interface{}); ok {
					location := geometry["location"].(map[string]interface{})
					latitude := location["lat"].(float64)
					longitude := location["lng"].(float64)

					// Update the database with the new coordinates
					db.Model(&scraper).Updates(Scraper{Latitude: fmt.Sprintf("%f", latitude), Longitude: fmt.Sprintf("%f", longitude)})
				}
			}
		}
	}
}
