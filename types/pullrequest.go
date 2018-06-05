package types

import (
	"fmt"
	"log"

	"github.com/ldez/go-git-cmd-wrapper/branch"
	"github.com/ldez/go-git-cmd-wrapper/checkout"
	"github.com/ldez/go-git-cmd-wrapper/fetch"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/pull"
	"github.com/ldez/go-git-cmd-wrapper/push"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/local"
	"github.com/pkg/errors"
)

// PullRequest the pull request model.
type PullRequest struct {
	Owner      string `json:"owner,omitempty"`
	BranchName string `json:"branch_name,omitempty"`
	Number     int    `json:"number,omitempty"`
	Project    string `json:"project,omitempty"`
	CloneURL   string `json:"clone_url,omitempty"`
}

const defaultInitialBranch = "master"

// Remove remove the pull request from the local git repository.
func (pr *PullRequest) Remove() error {
	// git remote get-url $remote
	out, err := git.Remote(remote.GetURL(pr.Owner))
	if err != nil {
		log.Println(out)
		return nil
	}

	branchName := makeLocalBranchName(pr)

	currentBranchName, err := local.GetCurrentBranchName()
	if err != nil {
		return errors.Wrapf(err, "[PR %d] unable to find current local branch name", pr.Number)
	}

	if currentBranchName == branchName {
		// git checkout $initial
		out, err = git.Checkout(checkout.Branch(defaultInitialBranch), git.Debug)
		if err != nil {
			log.Println(out)
			return errors.Wrapf(err, "[PR %d] unable to checkout initial branch (%s)", pr.Number, defaultInitialBranch)
		}
	}

	// git branch -D "$pr--$branch"
	out, err = git.Branch(branch.DeleteForce, branch.BranchName(branchName), git.Debug)
	if err != nil {
		log.Println(out)
		return errors.Wrapf(err, "[PR %d] unable to checkout PR", pr.Number)
	}

	return nil
}

// RemoveRemote remove the remote of the pull request from the local git repository.
func (pr *PullRequest) RemoveRemote() error {
	// git remote get-url $remote
	out, err := git.Remote(remote.GetURL(pr.Owner))
	if err != nil {
		log.Println(out)
		return nil
	}

	// git remote remove $remote
	out, err = git.Remote(remote.Remove(pr.Owner), git.Debug)
	if err != nil {
		log.Println(out)
		return errors.Wrapf(err, "[PR %d] unable to remove remote", pr.Number)
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
		return errors.Wrapf(err, "[PR %d] unable to push", pr.Number)
	}

	return nil
}

// Pull pull the PR from the remote git repository.
func (pr *PullRequest) Pull(force bool) error {
	// git pull -f $remote $branch
	out, err := git.Pull(git.Cond(force, pull.Force), pull.Repository(pr.Owner), pull.Refspec(pr.BranchName), git.Debug)
	if err != nil {
		log.Println(out)
		return errors.Wrapf(err, "[PR %d] unable to pull", pr.Number)
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
			if cloneURL == "" { // backward-compatible with previous configurations
				cloneURL = fmt.Sprintf("git@github.com:%s/%s.git", pr.Owner, pr.Project)
			}
			out, errRemote := git.Remote(remote.Add(pr.Owner, cloneURL), git.Debug)
			if errRemote != nil {
				log.Println(out)
				return errors.Wrapf(errRemote, "[PR %d] unable to add remote", pr.Number)
			}
		}

		// git fetch $remote $branch
		out, errFetch := git.Fetch(fetch.Remote(pr.Owner), fetch.RefSpec(pr.BranchName), git.Debug)
		if errFetch != nil {
			log.Println(out)
			return errors.Wrapf(errFetch, "[PR %d] unable to fetch", pr.Number)
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
		return errors.Wrapf(err, "[PR %d] unable to checkout", pr.Number)
	}

	return nil
}

func makeLocalBranchName(pr *PullRequest) string {
	return fmt.Sprintf("%d--%s", pr.Number, pr.BranchName)
}
