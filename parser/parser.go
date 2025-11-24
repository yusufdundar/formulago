/*
Copyright © 2022 Yusuf DÜNDAR <info@dundar.dev>

Parser for the formula1 official website
*/
package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode" // Added for name parsing

	"github.com/PuerkitoBio/goquery"
	"github.com/yusufdundar/formulago/model"
)

// fetchDocument performs an HTTP GET request for the given URL and returns
// a goquery.Document if successful.
// It includes a 30-second timeout for the request.
// Errors are returned if the request fails, the status code is not 200,
// or if HTML parsing fails.
func fetchDocument(url string) (*goquery.Document, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad status for URL %s: %s (status code %d)", url, res.Status, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from %s: %w", url, err)
	}
	return doc, nil
}

// parseName cleans a raw name string (typically a driver's or winner's name)
// by attempting to remove a trailing 3-letter uppercase code (e.g., "VER" from "Max VerstappenVER").
// If the trailing part is not 3 letters or not all uppercase, the original name is returned.
func parseName(rawName string) string {
	trimmedName := strings.TrimSpace(rawName)
	if len(trimmedName) > 3 {
		potentialCode := trimmedName[len(trimmedName)-3:]
		allUpper := true
		for _, r := range potentialCode {
			if !unicode.IsUpper(r) {
				allUpper = false
				break
			}
		}
		if allUpper {
			// Check if the character before the code is a letter, if not, it's likely a standalone code.
			// Example: "VER" should not become ""
			// Example: "Max VerstappenVER" should become "Max Verstappen"
			// Example: "Oscar PiastriPIA" should become "Oscar Piastri"
			if len(trimmedName)-4 >= 0 && (unicode.IsLetter(rune(trimmedName[len(trimmedName)-4])) || unicode.IsSpace(rune(trimmedName[len(trimmedName)-4]))) {
				return strings.TrimSpace(trimmedName[:len(trimmedName)-3])
			}
			// If it's something like "VER" alone and we don't want to strip it, we'd return trimmedName here.
			// However, the context is usually "FullNameCODE", so stripping is generally desired.
			// For now, if it's allUpper and 3 chars, we assume it's a code to be stripped from a longer name.
			// If the name itself is just "VER", this will make it empty. This might need adjustment
			// if driver codes themselves are ever primary identifiers without a name.
			// Given current usage, this mainly cleans names like "Max VerstappenVER".
			return strings.TrimSpace(trimmedName[:len(trimmedName)-3])
		}
	}
	return trimmedName // Return original trimmed name if no code found or name too short
}

// fetchLatestResultsYear fetches the HTML from initialUrl, parses it to find available result years,
// and returns the most recent year that is less than or equal to the current real-world year.
//
// The year selection logic primarily targets <a> tags within common filter component structures
// often found on websites for selecting seasons or years.
//
// Attempt 1 (Targeted): It first looks for <a> tags with href attributes containing "/en/results.html/"
// and validates if the link's text is a 4-digit year that also appears in the href.
// Example: <a href="/en/results.html/2024/drivers.html">2024</a>
//
// Attempt 2 (Broad Search - Fallback): If the targeted search yields no years, it iterates through ALL <a> tags.
// For each link, it logs its text and href for debugging purposes. It then validates if:
//  1. The link's text is a 4-digit number (e.g., between 1950 and current year + buffer for future).
//  2. The link's `href` attribute exists and contains "/results.html/" + the found year string + "/".
//
// This broad search helps identify potential year links even if their structure deviates significantly.
//
// The function sorts all unique, validated years in descending order and returns the first one
// that is less than or equal to the current real-world year (obtained via `time.Now().Year()`).
// If no such year is found, or if fetching/parsing fails, an error is returned.
func FetchLatestResultsYear(initialUrl string) (string, error) {
	// Temporarily override initialUrl for this specific test run
	fixedTestUrl := "https://www.formula1.com/en/results.html"

	doc, err := fetchDocument(fixedTestUrl)
	if err != nil {
		// Ensure the error message reflects the URL actually used for fetching
		return "", fmt.Errorf("error fetching document for year list from %s: %w", fixedTestUrl, err)
	}

	var years []int
	yearMap := make(map[int]bool) // To store unique years

	// Attempt 1: Targeted search for <a> tags with hrefs matching /en/results.html/YYYY/...
	// and text content being YYYY.
	doc.Find("a[href*='/en/results.html/']").Each(func(i int, s *goquery.Selection) {
		yearStr := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href") // href existence is guaranteed by the selector's attribute part

		if len(yearStr) == 4 {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				// Validate if the year from text is actually in the relevant part of the href
				if strings.Contains(href, "/results.html/"+yearStr+"/") {
					if !yearMap[year] {
						years = append(years, year)
						yearMap[year] = true
						// log.Printf("Targeted search: Found valid year %d from text '%s' and href '%s'", year, yearStr, href)
					}
				}
			}
		}
	})

	// Attempt 2: Broad search if targeted search yields no results or to augment (though yearMap handles duplicates)
	if len(years) == 0 {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			yearStrRaw := s.Text()
			href, exists := s.Attr("href")
			// Log the raw text and href

			yearStr := strings.TrimSpace(yearStrRaw)

			if len(yearStr) == 4 && exists { // Ensure href exists here
				year, err := strconv.Atoi(yearStr)
				if err == nil {
					// Basic sanity check for year range
					currentMaxYear := time.Now().Year() + 5 // Allow a small buffer for future years
					if year >= 1950 && year <= currentMaxYear {
						expectedPrefix := "/en/results/"           // Corrected prefix
						expectedYearSegment := "/" + yearStr + "/" // e.g., "/2024/"

						isValid := strings.HasPrefix(href, expectedPrefix) && strings.Contains(href, expectedYearSegment)

						if isValid {
							if !yearMap[year] {
								years = append(years, year)
								yearMap[year] = true
							}
						}
					}
				}

		})
	}

	if len(years) == 0 {
		// Ensure the error message reflects the URL actually used for fetching
		return "", errors.New("no years found in HTML from " + fixedTestUrl)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(years))) // Sort years in descending order

	currentRealYear := time.Now().Year()
	for _, year := range years {
		if year <= currentRealYear {
			return strconv.Itoa(year), nil
		}
	}
	// Ensure the error message reflects the URL actually used for fetching
	return "", fmt.Errorf("no valid year found (less than or equal to %d) from %s", currentRealYear, fixedTestUrl)
}

// ParseDriver Parse the driver standing info from formula1 website
func ParseDriver() []model.Driver {
	defaultYear := strconv.Itoa(time.Now().Year()) // Dynamic default year
	// Use a consistent page (like drivers) for fetching the latest year.
	yearFetchingURL := fmt.Sprintf("https://www.formula1.com/en/results.html/%s/drivers.html", defaultYear)

	latestYear, err := FetchLatestResultsYear(yearFetchingURL)
	if err != nil {
		log.Printf("Error fetching latest year for Drivers: %v. Falling back to default year %s.", err, defaultYear)
		latestYear = defaultYear
	}

	driverUrl := fmt.Sprintf("https://www.formula1.com/en/results.html/%s/drivers.html", latestYear)
	var DriverList []model.Driver

	doc, err := fetchDocument(driverUrl)
	if err != nil {
		log.Printf("Failed to fetch or parse driver data from %s: %v.", driverUrl, err)
		return DriverList // Return empty list on error
	}

	// New structure uses standard table
	data := doc.Find("table tbody tr")
	data.Each(func(i int, s *goquery.Selection) {
		pos := strings.TrimSpace(s.Find("td:nth-child(1)").Text())

		// Name is now in nested spans inside an anchor tag
		// Structure: <a ...><span><span class="max-lg:hidden">First</span> <span class="max-md:hidden">Last</span><span class="md:hidden">CODE</span></span></a>
		// We want First + Last.
		driverCell := s.Find("td:nth-child(2)")
		firstName := strings.TrimSpace(driverCell.Find("span.max-lg\\:hidden").Text())
		lastName := strings.TrimSpace(driverCell.Find("span.max-md\\:hidden").Text())
		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Fallback if specific classes aren't found (robustness)
		if firstName == "" && lastName == "" {
			name = strings.TrimSpace(driverCell.Text())
			// Clean up if it contains the code concatenated
			name = parseName(name)
		}

		nation := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
		car := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		pts := strings.TrimSpace(s.Find("td:nth-child(5)").Text())

		if pos != "" {
			driver := model.Driver{
				Pos:  pos,
				Name: name,
				Nat:  nation,
				Team: car,
				Pts:  pts,
			}
			DriverList = append(DriverList, driver)
		}
	})
	return DriverList
}

// ParseTeam Parse the constructor standing info from formula1 website
func ParseTeam() []model.Team {
	defaultYear := strconv.Itoa(time.Now().Year()) // Dynamic default year
	yearFetchingURL := fmt.Sprintf("https://www.formula1.com/en/results.html/%s/drivers.html", defaultYear)

	latestYear, err := FetchLatestResultsYear(yearFetchingURL)
	if err != nil {
		log.Printf("Error fetching latest year for Teams: %v. Falling back to default year %s.", err, defaultYear)
		latestYear = defaultYear
	}

	teamUrl := fmt.Sprintf("https://www.formula1.com/en/results.html/%s/team.html", latestYear)
	var TeamList []model.Team

	doc, err := fetchDocument(teamUrl)
	if err != nil {
		log.Printf("Failed to fetch or parse team data from %s: %v.", teamUrl, err)
		return TeamList // Return empty list on error
	}

	// Assumes Pos, Team Name, Pts are in columns 1, 2, 3 respectively.
	data := doc.Find("table tbody tr")
	data.Each(func(i int, s *goquery.Selection) {
		pos := strings.TrimSpace(s.Find("td:nth-child(1)").Text()) // Position in 1st column

		name := strings.TrimSpace(s.Find("td:nth-child(2) a").Text()) // Team Name in <a> tag in 2nd column
		if name == "" {
			// Fallback if team name is not in an <a> tag or <a> tag doesn't exist
			name = strings.TrimSpace(s.Find("td:nth-child(2)").Text())
		}

		pts := strings.TrimSpace(s.Find("td:nth-child(3)").Text()) // Points in 3rd column

		// Only add team if position is populated (simple check for valid row)
		if pos != "" {
			team := model.Team{
				Pos:  pos,
				Name: name,
				Pts:  pts,
			}
			TeamList = append(TeamList, team)
		}
	})
	return TeamList
}

// ParseRace Parse the F1 race info from formula1 website
func ParseRace() []model.Race {
	defaultYear := strconv.Itoa(time.Now().Year()) // Dynamic default year
	yearFetchingURL := fmt.Sprintf("https://www.formula1.com/en/results.html/%s/drivers.html", defaultYear)

	latestYear, err := FetchLatestResultsYear(yearFetchingURL)
	if err != nil {
		log.Printf("Error fetching latest year for Races: %v. Falling back to default year %s.", err, defaultYear)
		latestYear = defaultYear
	}

	raceUrl := fmt.Sprintf("https://www.formula1.com/en/results.html/%s/races.html", latestYear)
	var RaceList []model.Race

	doc, err := fetchDocument(raceUrl)
	if err != nil {
		log.Printf("Failed to fetch or parse race data from %s: %v.", raceUrl, err)
		return RaceList // Return empty list on error
	}

	// Assumes GrandPrix, Date, Winner, Car, Laps, Time are in columns 1-6.
	data := doc.Find("table tbody tr")
	data.Each(func(i int, s *goquery.Selection) {
		grandPrix := strings.TrimSpace(s.Find("td:nth-child(1) a").Text()) // Grand Prix from 1st column (usually a link)
		if grandPrix == "" {                                               // Fallback if not a link
			grandPrix = strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		}

		// Clean up race name by removing "flag of [country]" text
		// Example: "flag of Bahrain Bahrain Grand Prix" -> "Bahrain Grand Prix"
		if strings.HasPrefix(strings.ToLower(grandPrix), "flag of ") {
			// Remove "flag of [country] " prefix
			parts := strings.SplitN(grandPrix, " ", 4) // ["flag", "of", "Country", "Rest..."]
			if len(parts) >= 4 {
				grandPrix = strings.TrimSpace(parts[3])
			}
		}

		date := strings.TrimSpace(s.Find("td:nth-child(2)").Text()) // Date from 2nd column

		// Winner from 3rd column. Structure similar to driver name: nested spans in anchor.
		winnerCell := s.Find("td:nth-child(3)")
		firstName := strings.TrimSpace(winnerCell.Find("span.max-lg\\:hidden").Text())
		lastName := strings.TrimSpace(winnerCell.Find("span.max-md\\:hidden").Text())
		winner := fmt.Sprintf("%s %s", firstName, lastName)

		if firstName == "" && lastName == "" {
			winnerRaw := strings.TrimSpace(winnerCell.Text())
			winner = parseName(winnerRaw) // Use helper function as fallback
		}

		car := strings.TrimSpace(s.Find("td:nth-child(4)").Text())       // Car from 4th column
		laps := strings.TrimSpace(s.Find("td:nth-child(5)").Text())      // Laps from 5th column
		totalTime := strings.TrimSpace(s.Find("td:nth-child(6)").Text()) // Time from 6th column

		// Only add race if Grand Prix name is populated
		if grandPrix != "" {
			race := model.Race{
				GrandPrix: grandPrix,
				Date:      date,
				Winner:    winner,
				Car:       car,
				Laps:      laps,
				Time:      totalTime,
			}
			RaceList = append(RaceList, race)
		}
	})
	return RaceList
}
