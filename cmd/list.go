package cmd

import (
	"fmt"
	"os"

	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

func List(options *types.ListOptions) error {
	// get configuration
	configs, err := config.ReadFile()
	if err != nil {
		return err
	}

	if options.All {
		for _, conf := range configs {
			fmt.Println(conf.Directory)
			if len(conf.PullRequests) == 0 {
				fmt.Println("* 0 PR.")
			} else {
				for _, prs := range conf.PullRequests {
					for _, pr := range prs {
						fmt.Printf("* %d: %s - %s\n", pr.Number, pr.Owner, pr.BranchName)
					}
				}
			}
		}
	} else {
		repoDir, err := os.Getwd()
		if err != nil {
			return err
		}

		conf, err := config.Find(configs, repoDir)
		if err != nil {
			return err
		}

		if len(conf.PullRequests) == 0 {
			fmt.Println("* 0 PR.")
		} else {
			for _, prs := range conf.PullRequests {
				for _, pr := range prs {
					fmt.Printf("* %d: %s - %s\n", pr.Number, pr.Owner, pr.BranchName)
				}
			}
		}
	}

	return nil
}
