package main

import (
	"fmt"
	"log"

	"github.com/yusufdundar/formulago/parser"
)

func main() {
	testURL := "https://www.formula1.com/en/results.html/2024/drivers.html"
	log.Printf("Testing parser.FetchLatestResultsYear with URL: %s", testURL)

	year, err := parser.FetchLatestResultsYear(testURL)
	if err != nil {
		log.Fatalf("Error calling parser.FetchLatestResultsYear: %v", err)
	}
	fmt.Printf("Fetched latest year: %s\n", year)
}
