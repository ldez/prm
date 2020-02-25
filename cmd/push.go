package cmd

import (
	"context"
	"log"
	"strings"

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

	user, _, err := newGitHubClient(context.Background()).Users.Get(context.Background(), pr.Owner)
	if err != nil {
		return err
	}

	if strings.EqualFold(user.GetType(), "Organization") {
		log.Println("WARNING: GitHub has introduced a 'silent' breaking change:")
		log.Println("WARNING: it's now not possible to push on a fork's branch from a GitHub organization.")
	}

	return pr.Push(options.Force)
}
