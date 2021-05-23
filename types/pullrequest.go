package types

import (
	"fmt"
	"log"

	"github.com/ldez/go-git-cmd-wrapper/v2/branch"
	"github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	"github.com/ldez/go-git-cmd-wrapper/v2/fetch"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/pull"
	"github.com/ldez/go-git-cmd-wrapper/v2/push"
	"github.com/ldez/go-git-cmd-wrapper/v2/remote"
	"github.com/ldez/prm/v3/local"
)

// PullRequest the pull request model.
type PullRequest struct {
	Owner      string `json:"owner,omitempty"`
	BranchName string `json:"branch_name,omitempty"`
	Project    string `json:"project,omitempty"`
	CloneURL   string `json:"clone_url,omitempty"`
	Number     int    `json:"number,omitempty"`
}

const defaultInitialBranch = "master"

// Remove remove the pull request from the local git repository.
func (pr *PullRequest) Remove() error {
	// git remote get-url $remote
	out, err := git.Remote(remote.GetURL(pr.Owner))
	if err != nil {
		log.Println(out)
		// nolint:nilerr // ignore error
		return nil
	}

	branchName := makeLocalBranchName(pr)

	currentBranchName, err := local.GetCurrentBranchName()
	if err != nil {
		return fmt.Errorf("[PR %d] unable to find current local branch name: %w", pr.Number, err)
	}

	if currentBranchName == branchName {
		// git checkout $initial
		out, err = git.Checkout(checkout.Branch(defaultInitialBranch), git.Debug)
		if err != nil {
			log.Println(out)
			return fmt.Errorf("[PR %d] unable to checkout initial branch (%s): %w", pr.Number, defaultInitialBranch, err)
		}
	}

	// git branch -D "$pr--$branch"
	out, err = git.Branch(branch.DeleteForce, branch.BranchName(branchName), git.Debug)
	if err != nil {
		log.Println(out)
		return fmt.Errorf("[PR %d] unable to checkout PR: %w", pr.Number, err)
	}

	return nil
}

// RemoveRemote remove the remote of the pull request from the local git repository.
func (pr *PullRequest) RemoveRemote() error {
	// git remote get-url $remote
	out, err := git.Remote(remote.GetURL(pr.Owner))
	if err != nil {
		log.Println(out)
		// nolint:nilerr // ignore error
		return nil
	}

	// git remote remove $remote
	out, err = git.Remote(remote.Remove(pr.Owner), git.Debug)
	if err != nil {
		log.Println(out)
		return fmt.Errorf("[PR %d] unable to remove remote: %w", pr.Number, err)
	}

	return nil
}

// Push push the pull request to the remote git repository.
func (pr *PullRequest) Push(force bool) error {
	// git push --force-with-lease $remote $pr--$branch:$branch
	ref := fmt.Sprintf("%s:%s", makeLocalBranchName(pr), pr.BranchName)
	out, err := git.Push(push.NoFollowTags, git.Cond(force, push.ForceWithLease), push.Remote(pr.Owner), push.RefSpec(ref), git.Debug)
	if err != nil {
		log.Println(out)
		return fmt.Errorf("[PR %d] unable to push: %w", pr.Number, err)
	}

	return nil
}

// Pull pull the PR from the remote git repository.
func (pr *PullRequest) Pull(force bool) error {
	// git pull -f $remote $branch
	out, err := git.Pull(git.Cond(force, pull.Force), pull.Repository(pr.Owner), pull.Refspec(pr.BranchName), git.Debug)
	if err != nil {
		log.Println(out)
		return fmt.Errorf("[PR %d] unable to pull: %w", pr.Number, err)
	}

	return nil
}

// Checkout checkout the branch related to the pull request into the local git repository.
func (pr *PullRequest) Checkout(newBranch bool) error {
	if newBranch {
		// git remote get-url $remote
		_, err := git.Remote(remote.GetURL(pr.Owner))
		if err != nil {
			// git remote add $remote git@github.com:$remote/$project.git
			cloneURL := pr.CloneURL

			// backward-compatible with previous configurations
			if cloneURL == "" {
				cloneURL = fmt.Sprintf("git@github.com:%s/%s.git", pr.Owner, pr.Project)
			}

			out, errRemote := git.Remote(remote.Add(pr.Owner, cloneURL), git.Debug)
			if errRemote != nil {
				log.Println(out)
				return fmt.Errorf("[PR %d] unable to add remote: %w", pr.Number, errRemote)
			}
		}

		// git fetch $remote $branch
		out, errFetch := git.Fetch(fetch.Remote(pr.Owner), fetch.RefSpec(pr.BranchName), git.Debug)
		if errFetch != nil {
			log.Println(out)
			return fmt.Errorf("[PR %d] unable to fetch: %w", pr.Number, errFetch)
		}
	}

	// git checkout -t -b "$pr--$branch" $remote/$branch
	localBranchName := makeLocalBranchName(pr)
	startPoint := fmt.Sprintf("%s/%s", pr.Owner, pr.BranchName)
	out, err := git.Checkout(
		git.Cond(newBranch, checkout.Track, checkout.NewBranch(localBranchName), checkout.StartPoint(startPoint)),
		git.Cond(!newBranch, checkout.Branch(localBranchName)),
		git.Debug)
	if err != nil {
		log.Println(out)
		return fmt.Errorf("[PR %d] unable to checkout: %w", pr.Number, err)
	}

	return nil
}

func makeLocalBranchName(pr *PullRequest) string {
	return fmt.Sprintf("%d--%s", pr.Number, pr.BranchName)
}
