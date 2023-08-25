package main

import (
	"fmt"
	"time"
)

func main() {
		// Infinite loop to run the tasks every day at 11 PM
		for {
			// Calculate the next 11 PM
			now := time.Now()
			next11PM := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
			if next11PM.Before(now) {
				next11PM = next11PM.Add(24 * time.Hour)
			}
	
			// Calculate the duration until the next 11 PM
			durationUntilNext11PM := next11PM.Sub(now)
	
			// Sleep until the next 11 PM
			time.Sleep(durationUntilNext11PM)
	cbre()
	// Track the time for the arison() function
	// Always put arison first in the main function, if not it will mess up data.txt.
	fmt.Println("arison started")
	startTimeArison := time.Now()
	arison()
	endTimeArison := time.Now()
	totalTimeArison := endTimeArison.Sub(startTimeArison)
	fmt.Printf("arison execution time: %s\n", totalTimeArison)

	fmt.Println("cbre started")
	startTimeCBRE := time.Now()
	cbre()
	endTimeCBRE := time.Now()
	totalTimeCBRE := endTimeCBRE.Sub(startTimeCBRE)
	fmt.Printf("cbre execution time: %s\n", totalTimeCBRE)

	// Track the time for the loopnet() function
	fmt.Println("loopnet started")
	startTimeLoopnet := time.Now()
	loopnet()
	endTimeLoopnet := time.Now()
	totalTimeLoopnet := endTimeLoopnet.Sub(startTimeLoopnet)
	fmt.Printf("loopnet execution time: %s\n", totalTimeLoopnet)

	// Track the time for the omada() function
	fmt.Println("omada started")
	startTimeOmada := time.Now()
	omada()
	endTimeOmada := time.Now()
	totalTimeOmada := endTimeOmada.Sub(startTimeOmada)
	fmt.Printf("omada execution time: %s\n", totalTimeOmada)

	// Track the time for the royalpark() function
	fmt.Println("royalpark started")
	startTimeRoyalpark := time.Now()
	royalpark()
	endTimeRoyalpark := time.Now()
	totalTimeRoyalpark := endTimeRoyalpark.Sub(startTimeRoyalpark)
	fmt.Printf("royalpark execution time: %s\n", totalTimeRoyalpark)

	// Track the time for the cushman() function
	fmt.Println("cushman started")
	startTimeCushman := time.Now()
	cushman()
	endTimeCushman := time.Now()
	totalTimeCushman := endTimeCushman.Sub(startTimeCushman)
	fmt.Printf("cushman execution time: %s\n", totalTimeCushman)

	// Track the time for the ScrapeListingsFromMainURLs() function
	fmt.Println("spacelist started")
	startTimeScrape := time.Now()
	ScrapeListingsFromMainURLs()
	endTimeScrape := time.Now()
	totalTimeScrape := endTimeScrape.Sub(startTimeScrape)
	fmt.Printf("spacelist execution time: %s\n", totalTimeScrape)
			// Print a message to indicate the tasks have been executed
			fmt.Println("Tasks executed at", time.Now())

			// Sleep for the remaining duration of the day before starting the loop again
			durationUntilMidnight := time.Until(time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()))
			time.Sleep(durationUntilMidnight)
		}
}
