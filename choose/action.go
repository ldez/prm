package choose

import "github.com/AlecAivazis/survey/v2"

// Action name.
const (
	ActionList     = "list"
	ActionCheckout = "checkout"
	ActionRemove   = "remove"
	ActionProjects = "projects"
)

type answersAction struct {
	Action string
}

func (a answersAction) isExit() bool {
	return a.Action == ExitLabel
}

// Action get PRM action.
func Action() (string, error) {
	surveyOpts := []string{ActionList, ActionCheckout, ActionRemove, ActionProjects, ExitLabel}

	qs := []*survey.Question{
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "Choose the action",
				Options: surveyOpts,
				Help:    "https://ldez.github.io/prm/",
			},
		},
	}

	answers := &answersAction{}

	err := survey.Ask(qs, answers)
	if err != nil {
		return "", err
	}

	if answers.isExit() {
		return ExitLabel, nil
	}

	return answers.Action, nil
}
