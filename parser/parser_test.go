package parser

import (
	"testing"
)

func TestParseDriver(t *testing.T) {
	drivers := ParseDriver()
	if len(drivers) == 0 {
		t.Log("Warning: No drivers found. This might be expected if the season hasn't started or the site structure changed again.")
	} else {
		t.Logf("Found %d drivers", len(drivers))
		for i, d := range drivers {
			if i < 3 {
				t.Logf("Driver %d: %+v", i+1, d)
			}
			if d.Name == "" {
				t.Errorf("Driver at index %d has empty name", i)
			}
			if d.Pos == "" {
				t.Errorf("Driver at index %d has empty position", i)
			}
		}
	}
}

func TestParseTeam(t *testing.T) {
	teams := ParseTeam()
	if len(teams) == 0 {
		t.Log("Warning: No teams found.")
	} else {
		t.Logf("Found %d teams", len(teams))
		for i, team := range teams {
			if i < 3 {
				t.Logf("Team %d: %+v", i+1, team)
			}
			if team.Name == "" {
				t.Errorf("Team at index %d has empty name", i)
			}
			if team.Pos == "" {
				t.Errorf("Team at index %d has empty position", i)
			}
		}
	}
}

func TestParseRace(t *testing.T) {
	races := ParseRace()
	if len(races) == 0 {
		t.Log("Warning: No races found.")
	} else {
		t.Logf("Found %d races", len(races))
		for i, race := range races {
			if i < 3 {
				t.Logf("Race %d: %+v", i+1, race)
			}
			if race.GrandPrix == "" {
				t.Errorf("Race at index %d has empty Grand Prix", i)
			}
			if race.Winner == "" {
				t.Errorf("Race at index %d has empty Winner", i)
			}
		}
	}
}
