package choose

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	exitLabel = "exit"
	allLabel  = "all"
	// ExitValue representation
	ExitValue = 0
	// AllValue representation
	AllValue = math.MaxInt32
)

type answersPR struct {
	PR string
}

func (a answersPR) isExit() bool {
	return a.PR == exitLabel
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

type answersProject struct {
	Directory string
}

func (a answersProject) isExit() bool {
	return a.Directory == exitLabel
}

// PullRequest Choose a PR in the list
func PullRequest(pulls map[string][]types.PullRequest) (int, error) {
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
	surveyOpts = append(surveyOpts, allLabel)
	surveyOpts = append(surveyOpts, exitLabel)

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
		return ExitValue, nil
	}

	if answers.isAll() {
		return AllValue, nil
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

// RemotePulRequest Choose a PR in the list from GitHub
func RemotePulRequest(prs []*github.PullRequest) (int, error) {
	surveyOpts := []string{exitLabel}
	for _, pr := range prs {
		surveyOpts = append(surveyOpts, fmt.Sprintf("%d: %s - %s", pr.GetNumber(), pr.User.GetLogin(), pr.GetTitle()))
	}

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
		return ExitValue, nil
	}

	return answers.getPRNumber()
}
