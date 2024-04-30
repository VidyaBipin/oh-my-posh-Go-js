package segments

import (
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/mock"

	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
)

func TestITermSegment(t *testing.T) {
	cases := []struct {
		Case             string
		TermProgram      string
		SquelchMark      string
		Shell            string
		Template         string
		ExpectedString   string
		ExpectedDisabled bool
	}{
		{Case: "not iterm", TermProgram: "", SquelchMark: "1", Shell: "zsh", ExpectedDisabled: true},
		{Case: "default mark", TermProgram: "iTerm.app", Shell: "zsh", Template: "{{ .PromptMark }}", ExpectedDisabled: false},
		{Case: "zsh", TermProgram: "iTerm.app", SquelchMark: "1", Shell: "zsh", Template: "{{ .PromptMark }}", ExpectedString: `%{$(iterm2_prompt_mark)%}`},
		{Case: "bash", TermProgram: "iTerm.app", SquelchMark: "1", Shell: "bash", Template: "{{ .PromptMark }}", ExpectedString: `\[$(iterm2_prompt_mark)\]`},
		{Case: "fish", TermProgram: "iTerm.app", SquelchMark: "1", Shell: "fish", Template: "{{ .PromptMark }}", ExpectedDisabled: false},
		{Case: "pwsh", TermProgram: "iTerm.app", SquelchMark: "1", Shell: "pwsh", Template: "{{ .PromptMark }}", ExpectedDisabled: false},
		{Case: "gibberishshell", TermProgram: "iTerm.app", SquelchMark: "1", Shell: "jaserhuashf", Template: "{{ .PromptMark }}", ExpectedDisabled: false},
	}

	for _, tc := range cases {
		env := new(mock.MockedEnvironment)
		env.On("PathSeparator").Return("/")
		env.On("Getenv", "TERM_PROGRAM").Return(tc.TermProgram)
		env.On("Getenv", "ITERM2_SQUELCH_MARK").Return(tc.SquelchMark)
		env.On("Shell").Return(tc.Shell)
		env.On("Error", mock2.Anything).Return()
		iterm := &ITerm{
			env: env,
		}
		assert.Equal(t, !tc.ExpectedDisabled, iterm.Enabled(), tc.Case)
		if !tc.ExpectedDisabled {
			assert.Equal(t, tc.ExpectedString, renderTemplate(env, tc.Template, iterm), tc.Case)
		}
	}
}
