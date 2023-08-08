package main

import (
	"fmt"
	"time"
)

func main() {
	// Track the time for the arison() function
	fmt.Println("arison() started")
	startTimeArison := time.Now()
	arison()
	endTimeArison := time.Now()
	totalTimeArison := endTimeArison.Sub(startTimeArison)
	fmt.Printf("arison() execution time: %s\n", totalTimeArison)

	// Track the time for the ScrapeListingsFromMainURLs() function
	fmt.Println("ScrapeListingsFromMainURLs() started")
	startTimeScrape := time.Now()
	ScrapeListingsFromMainURLs()
	endTimeScrape := time.Now()
	totalTimeScrape := endTimeScrape.Sub(startTimeScrape)
	fmt.Printf("ScrapeListingsFromMainURLs() execution time: %s\n", totalTimeScrape)
}

