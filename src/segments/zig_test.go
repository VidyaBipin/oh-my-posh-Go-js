package segments

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/runtime"
	"github.com/stretchr/testify/assert"
)

func TestZig(t *testing.T) {
	cases := []struct {
		Case                  string
		Version               string
		InProjectDir          bool
		ExpectedString        string
		ExpectedMajor         string
		ExpectedMinor         string
		ExpectedPatch         string
		ExpectedPrerelease    string
		ExpectedBuildMetadata string
		ExpectedURL           string
	}{
		{
			Case:                  "zig 0.13.0 - not in project dir",
			Version:               "0.13.0",
			InProjectDir:          false,
			ExpectedString:        "0.13.0",
			ExpectedMajor:         "0",
			ExpectedMinor:         "13",
			ExpectedPatch:         "0",
			ExpectedPrerelease:    "",
			ExpectedBuildMetadata: "",
			ExpectedURL:           "https://ziglang.org/download/0.13.0/release-notes.html",
		},
		{
			Case:                  "zig 0.12.0-dev.2063+804cee3b9 - not in project dir",
			Version:               "0.12.0-dev.2063+804cee3b9",
			InProjectDir:          false,
			ExpectedString:        "0.12.0-dev.2063+804cee3b9",
			ExpectedMajor:         "0",
			ExpectedMinor:         "12",
			ExpectedPatch:         "0",
			ExpectedPrerelease:    "dev.2063",
			ExpectedBuildMetadata: "804cee3b9",
			ExpectedURL:           "https://ziglang.org/download/0.12.0/release-notes.html",
		},
		{
			Case:                  "zig 0.13.0 - in project dir",
			Version:               "0.13.0",
			InProjectDir:          true,
			ExpectedString:        "0.13.0",
			ExpectedMajor:         "0",
			ExpectedMinor:         "13",
			ExpectedPatch:         "0",
			ExpectedPrerelease:    "",
			ExpectedBuildMetadata: "",
			ExpectedURL:           "https://ziglang.org/download/0.13.0/release-notes.html",
		},
		{
			Case:                  "zig 0.12.0-dev.2063+804cee3b9 - in project dir",
			Version:               "0.12.0-dev.2063+804cee3b9",
			InProjectDir:          true,
			ExpectedString:        "0.12.0-dev.2063+804cee3b9",
			ExpectedMajor:         "0",
			ExpectedMinor:         "12",
			ExpectedPatch:         "0",
			ExpectedPrerelease:    "dev.2063",
			ExpectedBuildMetadata: "804cee3b9",
			ExpectedURL:           "https://ziglang.org/download/0.12.0/release-notes.html",
		},
	}
	for _, tc := range cases {
		params := &mockedLanguageParams{
			cmd:           "zig",
			versionParam:  "version",
			versionOutput: tc.Version,
			extension:     "*.zig",
		}
		env, props := getMockedLanguageEnv(params)

		dummyDir := &runtime.FileInfo{}

		if tc.InProjectDir {
			env.On("HasParentFilePath", "build.zig", false).Return(dummyDir, nil)
		} else {
			env.On("HasParentFilePath", "build.zig", false).Return(dummyDir, errors.New("build.zig not found"))
		}

		zig := &Zig{}
		zig.Init(props, env)

		failMsg := fmt.Sprintf("Failed in case: %s", tc.Case)
		assert.True(t, zig.Enabled(), failMsg)
		assert.Equal(t, tc.ExpectedString, renderTemplate(env, zig.Template(), zig), failMsg)
		assert.Equal(t, tc.ExpectedURL, renderTemplate(env, zig.URL, zig), failMsg)
		assert.Equal(t, tc.InProjectDir, zig.InProjectDir(), failMsg)
	}
}
