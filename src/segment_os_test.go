package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSInfo(t *testing.T) {
	cases := []struct {
		Case              string
		ExpectedString    string
		GOOS              string
		WSLDistro         string
		Platform          string
		DisplayDistroName bool
	}{
		{
			Case:           "WSL debian - icon",
			ExpectedString: "WSL at \uf306",
			GOOS:           "linux",
			WSLDistro:      "debian",
			Platform:       "debian",
		},
		{
			Case:              "WSL debian - name",
			ExpectedString:    "WSL at burps",
			GOOS:              "linux",
			WSLDistro:         "burps",
			Platform:          "debian",
			DisplayDistroName: true,
		},
		{
			Case:           "plain linux - icon",
			ExpectedString: "\uf306",
			GOOS:           "linux",
			Platform:       "debian",
		},
		{
			Case:              "plain linux - name",
			ExpectedString:    "debian",
			GOOS:              "linux",
			Platform:          "debian",
			DisplayDistroName: true,
		},
		{
			Case:           "windows",
			ExpectedString: "windows",
			GOOS:           "windows",
		},
		{
			Case:           "darwin",
			ExpectedString: "darwin",
			GOOS:           "darwin",
		},
		{
			Case:           "unknown",
			ExpectedString: "unknown",
			GOOS:           "unknown",
		},
	}
	for _, tc := range cases {
		env := new(MockedEnvironment)
		env.On("getRuntimeGOOS", nil).Return(tc.GOOS)
		env.On("getenv", "WSL_DISTRO_NAME").Return(tc.WSLDistro)
		env.On("getPlatform", nil).Return(tc.Platform)
		props := &properties{
			values: map[Property]interface{}{
				WSL:               "WSL",
				WSLSeparator:      " at ",
				DisplayDistroName: tc.DisplayDistroName,
				Windows:           "windows",
				MacOS:             "darwin",
			},
		}
		osInfo := &osInfo{
			env:   env,
			props: props,
		}
		assert.Equal(t, tc.ExpectedString, osInfo.string(), tc.Case)
		assert.Equal(t, tc.Platform, osInfo.OS, tc.Case)
	}
}
