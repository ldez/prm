package types

import (
	"fmt"
	"log"

	"github.com/ldez/go-git-cmd-wrapper/branch"
	"github.com/ldez/go-git-cmd-wrapper/checkout"
	"github.com/ldez/go-git-cmd-wrapper/fetch"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/push"
	"github.com/ldez/go-git-cmd-wrapper/remote"
)

// PullRequest the pull request model.
type PullRequest struct {
	Owner      string `json:"owner,omitempty"`
	BranchName string `json:"branch_name,omitempty"`
	Number     int    `json:"number,omitempty"`
	Project    string `json:"project,omitempty"`
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

	// git checkout $initial
	out, err = git.Checkout(checkout.Branch(defaultInitialBranch), git.Debug)
	if err != nil {
		log.Println(out)
		return err
	}

	// git branch -D "$pr--$branch"
	branchName := makeLocalBranchName(pr)
	out, err = git.Branch(branch.DeleteForce, branch.BranchName(branchName), git.Debug)
	if err != nil {
		log.Println(out)
		return err
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
		return err
	}

	return nil
}

// Push push force the pull request to the remote git repository.
func (pr *PullRequest) Push(force bool) error {

	// git push --force-with-lease $remote $pr--$branch:$branch
	ref := fmt.Sprintf("%s:%s", makeLocalBranchName(pr), pr.BranchName)
	out, err := git.Push(git.Cond(force, push.ForceWithLease), push.Remote(pr.Owner), push.RefSpec(ref), git.Debug)
	if err != nil {
		log.Println(out)
		return err
	}

	return nil
}

// Checkout checkout the branch related to the pull request into the local git repository.
func (pr *PullRequest) Checkout(newBranch bool) error {

	if newBranch {
		// git remote get-url $remote
		out, err := git.Remote(remote.GetURL(pr.Owner))
		if err != nil {
			// git remote add $remote git@github.com:$remote/$project.git
			forkURL := fmt.Sprintf("git@github.com:%s/%s.git", pr.Owner, pr.Project)
			out, err = git.Remote(remote.Add(pr.Owner, forkURL), git.Debug)
			if err != nil {
				log.Println(out)
				return err
			}
		}

		// git fetch $remote $branch
		out, err = git.Fetch(fetch.Remote(pr.Owner), fetch.RefSpec(pr.BranchName), git.Debug)
		if err != nil {
			log.Println(out)
			return err
		}
	}

	// git checkout -t -b "$pr--$branch" $remote/$branch
	localBranchName := makeLocalBranchName(pr)
	startPoint := fmt.Sprintf("%s/%s", pr.Owner, pr.BranchName)
	out, err := git.Checkout(
		git.Cond(newBranch, checkout.Track, checkout.NewBranch),
		checkout.Branch(localBranchName),
		git.Cond(newBranch, checkout.StartPoint(startPoint)),
		git.Debug)
	if err != nil {
		log.Println(out)
		return err
	}

	return nil
}

func makeLocalBranchName(pr *PullRequest) string {
	return fmt.Sprintf("%d--%s", pr.Number, pr.BranchName)
}
