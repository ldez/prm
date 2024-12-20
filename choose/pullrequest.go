package choose

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/go-github/v67/github"
	"github.com/ldez/prm/v3/types"
)

type answersPR struct {
	PR string
}

func (a answersPR) isExit() bool {
	return a.PR == ExitLabel
}

func (a answersPR) isAll() bool {
	return a.PR == allLabel
}

func (a answersPR) getPRNumber() (int, error) {
	parts := strings.SplitN(a.PR, ":", 2)

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	return number, nil
}

// PullRequest Choose a PR in the list.
func PullRequest(pulls map[string][]types.PullRequest, all bool) (int, error) {
	if len(pulls) == 0 {
		fmt.Println("* 0 PR.")
		return ExitValue, nil
	}

	var surveyOpts []string
	for _, prs := range pulls {
		for _, pr := range prs {
			surveyOpts = append(surveyOpts, fmt.Sprintf("%d: %s - %s", pr.Number, pr.Owner, pr.BranchName))
		}
	}
	sort.Strings(surveyOpts)
	if all {
		surveyOpts = append(surveyOpts, allLabel)
	}
	surveyOpts = append(surveyOpts, ExitLabel)

	qs := []*survey.Question{
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
		return ExitValue, nil
	}

	if answers.isAll() {
		return AllValue, nil
	}

	return answers.getPRNumber()
}

// RemotePulRequest Choose a PR in the list from GitHub.
func RemotePulRequest(prs []*github.PullRequest) (int, error) {
	surveyOpts := []string{ExitLabel}
	for _, pr := range prs {
		surveyOpts = append(surveyOpts, fmt.Sprintf("%d: %s - %s", pr.GetNumber(), pr.User.GetLogin(), pr.GetTitle()))
	}

	qs := []*survey.Question{
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
		return ExitValue, nil
	}

	return answers.getPRNumber()
}
