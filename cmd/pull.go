package cmd

import (
	"fmt"
	"os"

	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

// Pull pull the PR branch.
func Pull(options *types.PullOptions) error {

	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}

	con, err := config.Find(confs, repoDir)
	if err != nil {
		return err
	}

	number, err := getBranchPRNumber()
	if err != nil {
		return err
	}

	pr, err := con.FindPullRequests(number)
	if err != nil {
		return err
	}

	fmt.Println("pull", pr)

	err = pr.Pull(options.Force)
	if err != nil {
		return err
	}

	return nil
}
