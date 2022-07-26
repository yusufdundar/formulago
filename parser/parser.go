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

	"github.com/PuerkitoBio/goquery"
	"github.com/yusufdundar/formulago/model"
)

var driverUrl = "https://www.formula1.com/en/results.html/2022/drivers.html"
var teamUrl = "https://www.formula1.com/en/results.html/2022/team.html"
var raceUrl = "https://www.formula1.com/en/results.html/2022/races.html"

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
			data := doc.Find(".resultsarchive-table tbody tr:nth-child(-n+21)")

			data.Each(func(i int, s *goquery.Selection) {
				pos = s.Find("td:nth-child(2)").Text()
				s.Find("td:nth-child(3)").Each(func(j int, q *goquery.Selection) {
					name = q.Find("a > span.hide-for-mobile").Text()
				})
				nation = s.Find("td:nth-child(4)").Text()
				car = s.Find("td:nth-child(5) > a").Text()
				pts = s.Find("td:nth-child(6)").Text()

				driver := model.Driver{
					Pos:  pos,
					Name: name,
					Nat:  nation,
					Team: car,
					Pts:  pts,
				}
				DriverList = append(DriverList, driver)
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
			data := doc.Find(".resultsarchive-table tbody tr:nth-child(-n+10)")

			data.Each(func(i int, s *goquery.Selection) {
				pos = s.Find("td:nth-child(2)").Text()
				s.Find("td:nth-child(3)").Each(func(j int, q *goquery.Selection) {
					name = q.Find("a").Text()
				})
				pts = s.Find("td:nth-child(4)").Text()

				team := model.Team{
					Pos:  pos,
					Name: name,
					Pts:  pts,
				}
				TeamList = append(TeamList, team)
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
			data := doc.Find(".resultsarchive-table tbody tr:nth-child(-n+22)")

			data.Each(func(i int, s *goquery.Selection) {
				grandPrix = s.Find("td:nth-child(2) > a").Text()
				grandPrix = strings.TrimSpace(grandPrix)
				date = s.Find("td:nth-child(3)").Text()
				s.Find("td:nth-child(4)").Each(func(j int, q *goquery.Selection) {
					winner = q.Find("span.hide-for-mobile").Text()
				})
				car = s.Find("td:nth-child(5)").Text()
				laps = s.Find("td:nth-child(6)").Text()
				totalTime = s.Find("td:nth-child(7)").Text()

				race := model.Race{
					GrandPrix: grandPrix,
					Date:      date,
					Winner:    winner,
					Car:       car,
					Laps:      laps,
					Time:      totalTime,
				}
				RaceList = append(RaceList, race)
			})
		}
	}

	return RaceList
}
