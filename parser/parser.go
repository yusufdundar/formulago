/*
Copyright © 2022 Yusuf DÜNDAR <info@dundar.dev>

Parser for the formula1 official website
*/
package parser

import (
	"log"
	"net/http"
	"strings"
	"time"
	"unicode" // Added for name parsing

	"github.com/PuerkitoBio/goquery"
	"github.com/yusufdundar/formulago/model"
)

var driverUrl = "https://www.formula1.com/en/results.html/2024/drivers.html"
var teamUrl = "https://www.formula1.com/en/results.html/2024/team.html"
var raceUrl = "https://www.formula1.com/en/results.html/2024/races.html"

// ParseDriver Parse the driver standing info from formula1 website
func ParseDriver() []model.Driver {

	var pos = ""
	var name = ""
	var nation = ""
	var car = ""
	var pts = ""

	var DriverList []model.Driver

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(driverUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			// Adjusted selector to be more generic, removed row limit initially
			// Assumes the main table structure is similar, but column content/indices might change
			data := doc.Find(".resultsarchive-table tbody tr")

			data.Each(func(i int, s *goquery.Selection) {
				pos = strings.TrimSpace(s.Find("td:nth-child(1)").Text()) // Pos in 1st column

				// Driver Name parsing from 2nd column
				nameRaw := ""
				nameNode := s.Find("td:nth-child(2) a") // Name in <a> tag in 2nd column
				if nameNode.Length() > 0 {
					fullNameAndCode := strings.TrimSpace(nameNode.Text())
					if len(fullNameAndCode) > 3 {
						potentialCode := fullNameAndCode[len(fullNameAndCode)-3:]
						allUpper := true
						for _, r := range potentialCode {
							if !unicode.IsUpper(r) {
								allUpper = false
								break
							}
						}
						if allUpper {
							// Attempt to remove only the 3-letter code
							nameRaw = strings.TrimSpace(fullNameAndCode[:len(fullNameAndCode)-3])
						} else {
							nameRaw = fullNameAndCode // Fallback if not all upper (e.g. "De Vries")
						}
					} else {
						nameRaw = fullNameAndCode // Fallback for very short names
					}
				} else {
					// Fallback if <a> not found, try td text directly and clean it
					nameFromTd := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
					if len(nameFromTd) > 3 {
						potentialCode := nameFromTd[len(nameFromTd)-3:]
						allUpper := true
						for _, r := range potentialCode {
							if !unicode.IsUpper(r) {
								allUpper = false
								break
							}
						}
						if allUpper {
							nameRaw = strings.TrimSpace(nameFromTd[:len(nameFromTd)-3])
						} else {
							nameRaw = nameFromTd
						}
					} else {
						nameRaw = nameFromTd
					}
				}
				name = nameRaw

				nation = strings.TrimSpace(s.Find("td:nth-child(3)").Text()) // Nationality in 3rd column

				car = strings.TrimSpace(s.Find("td:nth-child(4) a").Text()) // Team in <a> tag in 4th column
				if car == "" {                                               // Fallback if team name is not a link
					car = strings.TrimSpace(s.Find("td:nth-child(4)").Text())
				}

				pts = strings.TrimSpace(s.Find("td:nth-child(5)").Text()) // Points in 5th column

				// Only add driver if position is populated (simple check for valid row)
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
		}
	}
	return DriverList
}

// ParseTeam Parse the constructor standing info from formula1 website
func ParseTeam() []model.Team {

	var pos = ""
	var name = ""
	var pts = ""

	var TeamList []model.Team

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(teamUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			// Adjusted selector: using the same table class, removed row limit for now.
			// Assumes Pos, Team Name, Pts are in columns 1, 2, 3 respectively.
			data := doc.Find(".resultsarchive-table tbody tr")

			data.Each(func(i int, s *goquery.Selection) {
				pos = strings.TrimSpace(s.Find("td:nth-child(1)").Text()) // Position in 1st column

				name = strings.TrimSpace(s.Find("td:nth-child(2) a").Text()) // Team Name in <a> tag in 2nd column
				if name == "" {
					// Fallback if team name is not in an <a> tag or <a> tag doesn't exist
					name = strings.TrimSpace(s.Find("td:nth-child(2)").Text())
				}

				pts = strings.TrimSpace(s.Find("td:nth-child(3)").Text()) // Points in 3rd column

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
		}

	}
	return TeamList
}

// ParseRace Parse the F1 race info from formula1 website
func ParseRace() []model.Race {

	var grandPrix = ""
	var date = ""
	var winner = ""
	var car = ""
	var laps = ""
	var totalTime = ""

	var RaceList []model.Race

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(raceUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			// Adjusted selector: using the same table class, removed row limit.
			// Assumes GrandPrix, Date, Winner, Car, Laps, Time are in columns 1-6.
			data := doc.Find(".resultsarchive-table tbody tr")

			data.Each(func(i int, s *goquery.Selection) {
				// Grand Prix from 1st column (usually a link)
				grandPrix = strings.TrimSpace(s.Find("td:nth-child(1) a").Text())
				if grandPrix == "" { // Fallback if not a link
					grandPrix = strings.TrimSpace(s.Find("td:nth-child(1)").Text())
				}

				date = strings.TrimSpace(s.Find("td:nth-child(2)").Text()) // Date from 2nd column

				// Winner from 3rd column, parse to remove 3-letter code
				winnerRaw := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
				if len(winnerRaw) > 3 {
					potentialCode := winnerRaw[len(winnerRaw)-3:]
					allUpper := true
					for _, r := range potentialCode {
						if !unicode.IsUpper(r) {
							allUpper = false
							break
						}
					}
					if allUpper {
						winner = strings.TrimSpace(winnerRaw[:len(winnerRaw)-3])
					} else {
						winner = winnerRaw // Fallback if code not all upper
					}
				} else {
					winner = winnerRaw // Fallback for short names
				}

				car = strings.TrimSpace(s.Find("td:nth-child(4)").Text())       // Car from 4th column
				laps = strings.TrimSpace(s.Find("td:nth-child(5)").Text())      // Laps from 5th column
				totalTime = strings.TrimSpace(s.Find("td:nth-child(6)").Text()) // Time from 6th column

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
		}
	}

	return RaceList
}
