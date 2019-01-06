package choose

import (
	"math"

	"github.com/ldez/prm/config"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const (
	allLabel = "all"
	// ExitLabel name
	ExitLabel = "exit"
	// ExitValue representation
	ExitValue = 0
	// AllValue representation
	AllValue = math.MaxInt32
)

type answersProject struct {
	Directory string
}

func (a answersProject) isExit() bool {
	return a.Directory == ExitLabel
}

// Project Choose a project in the list
func Project(configs []config.Configuration) (*config.Configuration, error) {
	var surveyOpts []string
	for _, conf := range configs {
		surveyOpts = append(surveyOpts, conf.Directory)
	}
	surveyOpts = append(surveyOpts, ExitLabel)

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
