package cmd

import (
	"fmt"

	"github.com/ldez/prm/config"
	"github.com/ldez/prm/local"
	"github.com/ldez/prm/types"
)

// List list PRs.
func List(options *types.ListOptions) error {
	// get configuration
	configs, err := config.ReadFile()
	if err != nil {
		return err
	}

	if options.All {
		displayProjects(configs)
	} else {
		repoDir, err := local.GetGitRepoRoot()
		if err != nil {
			return err
		}

		conf, err := config.Find(configs, repoDir)
		if err != nil {
			return err
		}

		displayPullRequests(conf.PullRequests)
	}

	return nil
}

func displayPullRequests(pulls map[string][]types.PullRequest) {
	if len(pulls) == 0 {
		fmt.Println("* 0 PR.")
	} else {
		for _, prs := range pulls {
			for _, pr := range prs {
				fmt.Printf("* %d: %s - %s\n", pr.Number, pr.Owner, pr.BranchName)
			}
		}
	}
}

func displayProjects(configs []config.Configuration) {
	for _, conf := range configs {
		fmt.Println(conf.Directory)
		displayPullRequests(conf.PullRequests)
	}
}
