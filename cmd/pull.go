package cmd

import (
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/local"
	"github.com/ldez/prm/types"
)

// Pull pull the PR branch.
func Pull(options *types.PullOptions) error {
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

	number, err := local.GetCurrentBranchPRNumber()
	if err != nil {
		return err
	}

	pr, err := con.FindPullRequests(number)
	if err != nil {
		return err
	}

	return pr.Pull(options.Force)
}
