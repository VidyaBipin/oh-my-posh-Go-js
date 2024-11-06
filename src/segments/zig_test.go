package segments

import (
	"errors"
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

		assert.True(t, zig.Enabled(), tc.Case)
		assert.Equal(t, tc.ExpectedString, renderTemplate(env, zig.Template(), zig), tc.Case)
		assert.Equal(t, tc.ExpectedMajor, renderTemplate(env, zig.Major, zig), tc.Case)
		assert.Equal(t, tc.ExpectedMinor, renderTemplate(env, zig.Minor, zig), tc.Case)
		assert.Equal(t, tc.ExpectedPatch, renderTemplate(env, zig.Patch, zig), tc.Case)
		assert.Equal(t, tc.ExpectedPrerelease, renderTemplate(env, zig.Prerelease, zig), tc.Case)
		assert.Equal(t, tc.ExpectedBuildMetadata, renderTemplate(env, zig.BuildMetadata, zig), tc.Case)
		assert.Equal(t, tc.ExpectedURL, renderTemplate(env, zig.URL, zig), tc.Case)
		assert.Equal(t, tc.InProjectDir, zig.InProjectDir(), tc.Case)
	}
}
