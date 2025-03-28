package choose

import "github.com/AlecAivazis/survey/v2"

// Fork Choose to create a fork automatically.
func Fork() (bool, error) {
	prompt := &survey.Confirm{
		Message: "Do you want PRM to create a fork automatically?",
	}

	yes := false

	err := survey.AskOne(prompt, &yes)
	if err != nil {
		return false, err
	}

	return yes, nil
}

// Username Asking for the GitHub username.
func Username() (string, error) {
	prompt := &survey.Input{
		Message: "Your GitHub username",
	}

	username := ""
	err := survey.AskOne(prompt, &username)

	return username, err
}
