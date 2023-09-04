package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		currentTime := time.Now()
		nextExecutionTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 18, 0, 0, 0, currentTime.Location())
		if currentTime.After(nextExecutionTime) {
			nextExecutionTime = nextExecutionTime.Add(24 * time.Hour)
		}

		sleepDuration := nextExecutionTime.Sub(currentTime)
		fmt.Printf("Next execution time: %s\n", nextExecutionTime)
		fmt.Printf("Sleeping for %s\n", sleepDuration)

		time.Sleep(sleepDuration)
	// Track the time for the arison() function
	// Always put arison first in the main function, if not it will mess up data.txt.
	fmt.Println("arison started")
	startTimeArison := time.Now()
	arison()
	endTimeArison := time.Now()
	totalTimeArison := endTimeArison.Sub(startTimeArison)
	fmt.Printf("arison execution time: %s\n", totalTimeArison)

	fmt.Println("yeg started")
	startTimeYeg := time.Now()
	yeg()
	endTimeYeg := time.Now()
	totalTimeYeg := endTimeYeg.Sub(startTimeYeg)
	fmt.Printf("arison execution time: %s\n", totalTimeYeg)
		
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

}
		}
