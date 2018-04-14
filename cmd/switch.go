package cmd

import (
	"os"
	"strings"

	"github.com/ldez/prm/choose"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

// Switch from a list.
func Switch(options *types.ListOptions) error {
	// get configuration
	configs, err := config.ReadFile()
	if err != nil {
		return err
	}

	if options.All {
		err := changeProject(configs)
		if err != nil {
			return err
		}
	} else {
		err := changePR(configs)
		if err != nil {
			return err
		}
	}

	return nil
}

func changeProject(configs []config.Configuration) error {
	var projectsConf []config.Configuration
	for _, value := range configs {
		if len(value.PullRequests) > 0 {
			projectsConf = append(projectsConf, value)
		}
	}

	conf, err := choose.Project(projectsConf)
	if err != nil || conf == nil {
		return err
	}

	number, err := choose.PullRequest(conf.PullRequests, false)
	if err != nil || number <= choose.ExitValue {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(currentDir, conf.Directory) {
		err := os.Chdir(conf.Directory)
		if err != nil {
			return err
		}
	}
	return Checkout(&types.CheckoutOptions{Number: number})
}

func changePR(configs []config.Configuration) error {
	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}

	conf, err := config.Find(configs, repoDir)
	if err != nil {
		return err
	}

	number, err := choose.PullRequest(conf.PullRequests, false)
	if err != nil || number == choose.ExitValue {
		return err
	}
	return Checkout(&types.CheckoutOptions{Number: number})
}
