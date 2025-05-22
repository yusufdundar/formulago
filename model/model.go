/*
Copyright © 2022 Yusuf DÜNDAR <info@dundar.dev>

*/
package model

// Driver represents a Formula 1 driver's standings information.
type Driver struct {
	Pos  string // Position
	Name string // Driver's full name
	Nat  string // Nationality
	Team string // Team/Constructor
	Pts  string // Points
}

// Team represents a Formula 1 constructor's standings information.
type Team struct {
	Pos  string // Position
	Name string // Team/Constructor name
	Pts  string // Points
}

// Race represents information about a specific Formula 1 race result.
type Race struct {
	GrandPrix string // Name of the Grand Prix
	Date      string // Date of the race
	Winner    string // Full name of the winning driver
	Car       string // Winning constructor/car
	Laps      string // Number of laps completed
	Time      string // Total race time of the winner
}
