package segments

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
	testify_ "github.com/stretchr/testify/mock"
)

func getTestData(file string) string {
	content, _ := os.ReadFile(fmt.Sprintf("../test/%s", file))
	return string(content)
}

// create Test segment for NBA segment
func TestNBASegment(t *testing.T) {
	jsonScheduleData := getTestData("nba/schedule.json")
	jsonScoreData := getTestData("nba/score.json")

	cases := []struct {
		Error           error
		Case            string
		JSONResponse    string
		ExpectedString  string
		TeamName        string
		CacheTimeout    int
		DaysOffset      int
		ExpectedEnabled bool
		CacheFoundFail  bool
	}{
		{
			Case:            "Team (Home Team) Scheduled Game",
			JSONResponse:    jsonScheduleData,
			TeamName:        "LAL",
			ExpectedString:  "󰠆 LAL vs PHX | 10/26/2023 | 10:00 PM ET",
			ExpectedEnabled: true,
			DaysOffset:      8,
		},
		{
			Case:            "Team (Away Team) Scheduled Game",
			JSONResponse:    jsonScheduleData,
			TeamName:        "PHX",
			ExpectedString:  "󰠆 LAL vs PHX | 10/26/2023 | 10:00 PM ET",
			DaysOffset:      4,
			ExpectedEnabled: true,
		},
		{
			Case:            "Team (Home Team) Live Game",
			JSONResponse:    jsonScoreData,
			TeamName:        "CHA",
			ExpectedString:  "󰠆 CHA (1-0):13 vs BOS (0-1):8 | Q1 8:23",
			ExpectedEnabled: true,
		},
		{
			Case:            "Team (Away Team) Live Game",
			JSONResponse:    jsonScoreData,
			TeamName:        "BOS",
			ExpectedString:  "󰠆 CHA (1-0):13 vs BOS (0-1):8 | Q1 8:23",
			ExpectedEnabled: true,
		},
		{
			Case:            "Team not Found",
			JSONResponse:    jsonScheduleData,
			DaysOffset:      8,
			TeamName:        "INVALID",
			ExpectedEnabled: false,
		},
	}

	for _, tc := range cases {
		env := &mock.Environment{}
		props := properties.Map{
			TeamName:   tc.TeamName,
			DaysOffset: tc.DaysOffset,
		}

		env.On("Error", testify_.Anything)
		env.On("Debug", testify_.Anything)
		env.On("HTTPRequest", NBAScoreURL).Return([]byte(tc.JSONResponse), tc.Error)

		// Add all the daysOffset to the http request responses
		for i := 0; i < tc.DaysOffset; i++ {
			currTime := time.Now().In(time.FixedZone("America/New_York", -5*60*60))
			// add offset days to currTime so we can query for games in the future
			currTime = currTime.AddDate(0, 0, i)
			dateStr := currTime.Format(NBADateFormat)
			scheduleURLEndpoint := fmt.Sprintf(NBAScheduleURL, currentNBASeason, dateStr)
			env.On("HTTPRequest", scheduleURLEndpoint).Return([]byte(tc.JSONResponse), tc.Error)
		}

		nba := &Nba{}
		nba.Init(props, env)

		enabled := nba.Enabled()
		assert.Equal(t, tc.ExpectedEnabled, enabled, tc.Case)
		if !enabled {
			continue
		}

		assert.Equal(t, tc.ExpectedString, renderTemplate(env, nba.Template(), nba), tc.Case)
	}
}
