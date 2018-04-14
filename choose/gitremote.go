package choose

import (
	"fmt"
	"strings"

	"github.com/ldez/prm/local"
	"gopkg.in/AlecAivazis/survey.v1"
)

type answersGitRemote struct {
	Remote string
}

func (a answersGitRemote) isExit() bool {
	return a.Remote == ExitLabel
}

func (a answersGitRemote) getName() string {
	parts := strings.SplitN(a.Remote, "]:", 2)
	return strings.TrimPrefix(parts[0], "[")
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
				Message: "Choose the remote related to Pull Requests (main remote)",
				Options: surveyOpts,
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

	return answers.getName(), nil
}
