package cmd

import (
	"github.com/ldez/prm/v3/config"
	"github.com/ldez/prm/v3/local"
	"github.com/ldez/prm/v3/types"
)

// Push push to the PR branch.
func Push(options *types.PushOptions) error {
	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := local.GetGitRepoRoot()
	if err != nil {
		return err
	}

	con, err := config.Find(confs, repoDir)
	if err != nil {
		return err
	}

	number, err := local.GetCurrentPRNumber(options.Number)
	if err != nil {
		return err
	}

	pr, err := con.FindPullRequests(number)
	if err != nil {
		return err
	}

	return pr.Push(options.Force)
}
