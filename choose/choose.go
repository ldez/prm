package choose

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

type answersPR struct {
	PR string
}

func (a answersPR) isExit() bool {
	return a.PR == "exit"
}

func (a answersPR) getPRNumber() (int, error) {
	parts := strings.SplitN(a.PR, ":", 2)

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	return number, nil
}

type answersProject struct {
	Directory string
}

func (a answersProject) isExit() bool {
	return a.Directory == "exit"
}

// PullRequest Choose a PR in the list
func PullRequest(pulls map[string][]types.PullRequest) (int, error) {
	if len(pulls) == 0 {
		fmt.Println("* 0 PR.")
		return 0, nil
	}

	var surveyOpts []string
	for _, prs := range pulls {
		for _, pr := range prs {
			surveyOpts = append(surveyOpts, fmt.Sprintf("%d: %s - %s", pr.Number, pr.Owner, pr.BranchName))
		}
	}
	surveyOpts = append(surveyOpts, "exit")

	var qs = []*survey.Question{
		{
			Name: "pr",
			Prompt: &survey.Select{
				Message: "Choose a PR",
				Options: surveyOpts,
			},
		},
	}

	answers := &answersPR{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return 0, err
	}
	if answers.isExit() {
		return 0, nil
	}

	return answers.getPRNumber()
}

// Project Choose a project in the list
func Project(configs []config.Configuration) (*config.Configuration, error) {
	var surveyOpts []string
	for _, conf := range configs {
		surveyOpts = append(surveyOpts, conf.Directory)
	}
	surveyOpts = append(surveyOpts, "exit")

	var qs = []*survey.Question{
		{
			Name: "directory",
			Prompt: &survey.Select{
				Message: "Choose a directory",
				Options: surveyOpts,
			},
		},
	}

	answers := &answersProject{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return nil, err
	}

	if answers.isExit() {
		return nil, nil
	}
	return config.Find(configs, answers.Directory)
}
