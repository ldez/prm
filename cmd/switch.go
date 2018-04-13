package cmd

import (
	"math"
	"os"
	"strings"

	"github.com/ldez/prm/choose"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

// Switch from the list of PRs.
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
	conf, err := choose.Project(configs)
	if err != nil || conf == nil {
		return err
	}

	number, err := choose.PullRequest(conf.PullRequests)
	if err != nil || number <= 0 {
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

	number, err := choose.PullRequest(conf.PullRequests)
	if err != nil || number <= 0 || number == math.MaxInt32 {
		return err
	}
	return Checkout(&types.CheckoutOptions{Number: number})
}
