package choose

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ldez/prm/v3/local"
)

type answersGitRemote struct {
	Remote string
}

func (a answersGitRemote) isExit() bool {
	return a.Remote == ExitLabel
}

// GitRemote Choose the remote related to PRs (main remote)
func GitRemote(remotes []local.Remote) (string, error) {
	var surveyOpts []string
	for _, remote := range remotes {
		surveyOpts = append(surveyOpts, fmt.Sprintf("[%s]: %s", remote.Name, remote.URL))
	}
	surveyOpts = append(surveyOpts, ExitLabel)

	var qs = []*survey.Question{
		{
			Name: "remote",
			Prompt: &survey.Select{
				Message: "Choose the remote",
				Options: surveyOpts,
				Help:    "The remote must be the repository where are the Pull Requests.",
			},
			Transform: func(ans interface{}) interface{} {
				answer, ok := ans.(survey.OptionAnswer)
				if !ok {
					return nil
				}
				parts := strings.SplitN(answer.Value, "]:", 2)
				answer.Value = strings.TrimPrefix(parts[0], "[")
				return answer
			},
		},
	}

	answers := &answersGitRemote{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return "", err
	}

	if answers.isExit() {
		return ExitLabel, nil
	}

	return answers.Remote, nil
}
