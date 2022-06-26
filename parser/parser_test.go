/*
Copyright © 2022 Yusuf DÜNDAR <info@dundar.dev>

*/
package parser

import "testing"

func TestParseDriver(t *testing.T) {

	got := ParseDriver()
	want := 20

	if len(got) <= want {
		t.Errorf("got %q, wanted %q", len(got), want)
	}
}

func TestParseTeam(t *testing.T) {

	got := ParseTeam()
	want := 10

	if len(got) != want {
		t.Errorf("got %q, wanted %q", len(got), want)
	}
}

func TestParseRace(t *testing.T) {

	got := ParseRace()
	notWant := 0

	if len(got) == notWant {
		t.Errorf("got %q, but not wanted %q", len(got), notWant)
	}
}
