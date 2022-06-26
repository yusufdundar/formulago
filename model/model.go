/*
Copyright © 2022 Yusuf DÜNDAR <info@dundar.dev>

*/
package model

type Driver struct {
	Pos  string
	Name string
	Nat  string
	Team string
	Pts  string
}

type Team struct {
	Pos  string
	Name string
	Pts  string
}

type Race struct {
	GrandPrix string
	Date      string
	Winner    string
	Car       string
	Laps      string
	Time      string
}
