package choose

import (
	"github.com/AlecAivazis/survey/v2"
)

type answersMainBranch struct {
	Branch string
}

func (a answersMainBranch) isExit() bool {
	return a.Branch == ExitLabel
}

// MainBranch Choose the main branch of the repository.
func MainBranch(branches []string) (string, error) {
	var surveyOpts []string
	surveyOpts = append(surveyOpts, branches...)
	surveyOpts = append(surveyOpts, ExitLabel)

	qs := []*survey.Question{
		{
			Name: "branch",
			Prompt: &survey.Select{
				Message: "Choose the main branch",
				Options: surveyOpts,
				Help:    "The main branch is the default branch ('main', 'master', ...).",
			},
		},
	}

	answers := &answersMainBranch{}
	err := survey.Ask(qs, answers)
	if err != nil {
		return "", err
	}

	if answers.isExit() {
		return ExitLabel, nil
	}

	return answers.Branch, nil
}
