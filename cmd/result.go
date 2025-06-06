/*
Copyright © 2022 Yusuf DÜNDAR <info@dundar.dev>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yusufdundar/formulago/parser"
)

const DATE_FORMAT_STRING = "02.01.2006 15:04:05"

var Driver bool
var Constructor bool
var Race bool

// resultCmd represents the result command
var resultCmd = &cobra.Command{
	Use:   "result",
	Short: "Displays the standings of drivers or constructors and race result",
	Long: `This command displays current statistics about F1 races 
	such as driver, constructor or race in your terminal.`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if Driver {
			showDriver()
		} else if Constructor {
			showConstructor()
		} else if Race {
			showRace()
		} else {
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func showDriver() {
	captionText := "\nCurrent Driver standings as of " + time.Now().Format(DATE_FORMAT_STRING)
	fmt.Println(captionText)

	// Fetch driver data using the parser
	driverData := parser.ParseDriver()

	// Check if data was retrieved
	if len(driverData) == 0 {
		fmt.Println("No driver data found or an error occurred while fetching.")
		return
	}

	// Prepare data for table
	var data [][]string
	for _, driver := range driverData {
		data = append(data, []string{driver.Pos, driver.Name, driver.Nat, driver.Team, driver.Pts})
	}

	// Create and render table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"POS", "DRIVER", "NATION", "TEAM", "PTS"})
	table.AppendBulk(data)
	table.Render()
}

func showConstructor() {
	captionText := "\nCurrent constructor standings as of " + time.Now().Format(DATE_FORMAT_STRING)
	fmt.Println(captionText)

	// Fetch team data using the parser
	teamData := parser.ParseTeam()

	// Check if data was retrieved
	if len(teamData) == 0 {
		fmt.Println("No constructor data found or an error occurred while fetching.")
		return
	}

	// Prepare data for table
	var data [][]string
	for _, team := range teamData {
		data = append(data, []string{team.Pos, team.Name, team.Pts})
	}

	// Create and render table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"POS", "TEAM", "PTS"})
	table.AppendBulk(data)
	table.Render()
}

func showRace() {
	captionText := "\nRace results as of " + time.Now().Format(DATE_FORMAT_STRING)
	fmt.Println(captionText)

	// Fetch race data using the parser
	raceData := parser.ParseRace()

	// Check if data was retrieved
	if len(raceData) == 0 {
		fmt.Println("No race data found or an error occurred while fetching.")
		return
	}

	// Prepare data for table
	var data [][]string
	for _, race := range raceData {
		data = append(data, []string{race.GrandPrix, race.Date, race.Winner, race.Car, race.Laps, race.Time})
	}

	// Create and render table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"GRAND PRIX", "DATE", "WINNER", "CAR", "LAPS", "TIME"})
	table.AppendBulk(data)
	table.Render()
}

func init() {
	rootCmd.AddCommand(resultCmd)

	resultCmd.Flags().BoolVarP(&Driver, "driver", "d", false, "Display driver standings")
	resultCmd.Flags().BoolVarP(&Constructor, "constructor", "c", false, "Display constructor standings")
	resultCmd.Flags().BoolVarP(&Race, "race", "r", false, "Display race results")

}
