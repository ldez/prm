package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/go-github/v30/github"
	"github.com/ldez/prm/v3/choose"
	"github.com/ldez/prm/v3/config"
	"github.com/ldez/prm/v3/local"
	"github.com/ldez/prm/v3/types"
)

// InteractiveCheckout checkout a PR.
func InteractiveCheckout(conf *config.Configuration) error {
	baseRepository, err := types.GetRepository(conf.BaseRemote)
	if err != nil {
		return err
	}

	// Display PRs from GitHub
	number, err := getPRNumberFromGitHub(baseRepository)
	if err != nil {
		return err
	}

	if number == choose.ExitValue {
		return nil
	}

	checkoutOptions := &types.CheckoutOptions{
		Number: number,
	}
	return Checkout(checkoutOptions)
}

func getPRNumberFromGitHub(baseRepository *types.Repository) (int, error) {
	ctx := context.Background()
	client := newGitHubClient(ctx)

	opt := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 50},
	}

	prs, _, err := client.PullRequests.List(ctx, baseRepository.Owner, baseRepository.Name, opt)
	if err != nil {
		return 0, fmt.Errorf("fail to retrieve pull request from GitHub: %w", err)
	}

	return choose.RemotePulRequest(prs)
}

// Checkout checkout a PR.
func Checkout(options *types.CheckoutOptions) error {
	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := local.GetGitRepoRoot()
	if err != nil {
		return err
	}

	conf, err := config.Find(confs, repoDir)
	if err != nil {
		return err
	}

	// check if already exists in config
	pr, err := conf.FindPullRequests(options.Number)
	if err == nil {
		log.Println("PR already exists.")

		// simple checkout
		return pr.Checkout(false)
	}

	log.Println("New Pull Request.")

	baseRepository, err := types.GetRepository(conf.BaseRemote)
	if err != nil {
		return err
	}

	pr, err = getPullRequest(baseRepository, options.Number)
	if err != nil {
		return err
	}

	err = pr.Checkout(true)
	if err != nil {
		// Remove remote if needed
		errRemote := removeRemote(conf, pr)
		if errRemote != nil {
			log.Println(errRemote)
		}
		return err
	}

	// add PR to config
	if conf.PullRequests == nil {
		conf.PullRequests = make(map[string][]types.PullRequest)
	}
	conf.PullRequests[pr.Owner] = append(conf.PullRequests[pr.Owner], *pr)

	return config.Save(confs)
}

// removeRemote if needed
func removeRemote(conf *config.Configuration, pr *types.PullRequest) error {
	if len(conf.PullRequests[pr.Owner]) == 0 {
		errRemote := pr.RemoveRemote()
		if errRemote != nil {
			return errRemote
		}
	}
	return nil
}

func getPullRequest(baseRepository *types.Repository, number int) (*types.PullRequest, error) {
	ctx := context.Background()
	client := newGitHubClient(ctx)

	pr, _, err := client.PullRequests.Get(ctx, baseRepository.Owner, baseRepository.Name, number)
	if err != nil {
		return nil, err
	}

	if pr.Head == nil || pr.Head.Repo == nil || pr.Head.Repo.Owner == nil {
		return nil, errors.New("the repository of the pull request has been deleted")
	}

	return &types.PullRequest{
		Project:    baseRepository.Name,
		Owner:      pr.Head.Repo.Owner.GetLogin(),
		BranchName: pr.Head.GetRef(),
		Number:     number,
		CloneURL:   pr.Head.Repo.GetSSHURL(),
	}, nil
}
