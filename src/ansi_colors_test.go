package main

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestGetAnsiFromColorString(t *testing.T) {
	cases := []struct {
		Case       string
		Expected   AnsiColor
		Color      string
		Background bool
	}{
		{Case: "Invalid background", Expected: emptyAnsiColor, Color: "invalid", Background: true},
		{Case: "Invalid background", Expected: emptyAnsiColor, Color: "invalid", Background: false},
		{Case: "Hex foreground", Expected: AnsiColor("38;2;170;187;204"), Color: "#AABBCC", Background: false},
		{Case: "Hex backgrond", Expected: AnsiColor("48;2;170;187;204"), Color: "#AABBCC", Background: true},
		{Case: "Base 8 foreground", Expected: AnsiColor("31"), Color: "red", Background: false},
		{Case: "Base 8 background", Expected: AnsiColor("41"), Color: "red", Background: true},
		{Case: "Base 16 foreground", Expected: AnsiColor("91"), Color: "lightRed", Background: false},
		{Case: "Base 16 backround", Expected: AnsiColor("101"), Color: "lightRed", Background: true},
	}
	for _, tc := range cases {
		ansiColors := &DefaultAnsiColors{}
		ansiColor := ansiColors.AnsiColorFromString(tc.Color, tc.Background)
		assert.Equal(t, tc.Expected, ansiColor, tc.Case)
	}
}
