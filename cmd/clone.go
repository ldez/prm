package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/fetch"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/v3/choose"
	"github.com/ldez/prm/v3/types"
	giturls "github.com/whilp/git-urls"
)

// Clone clone and fork a repository.
func Clone(options types.CloneOptions) error {
	srcRepoURL, err := giturls.Parse(options.Repo)
	if err != nil {
		return err
	}

	user, repoName, err := splitUserRepo(srcRepoURL)
	if err != nil {
		return err
	}

	if options.UserAsRootDir {
		if err = os.MkdirAll(filepath.Clean(user), 0755); err != nil {
			return err
		}

		if err = os.Chdir(filepath.Clean(user)); err != nil {
			return err
		}
	}

	if options.NoFork {
		return simpleClone(options)
	}

	if !HasToken() {
		fmt.Println("---------------------------------------------------------")
		fmt.Printf("Set %s or %s to allow to fork automatically:\n", tokenEnvVar, tokenEnvVar+fileSuffixEnvVar)
		fmt.Println("- https://ldez.github.io/prm/#prm-github-token")
		fmt.Println("- https://ldez.github.io/prm/#prm-github-token-file")
		fmt.Println("---------------------------------------------------------")
		return simpleClone(options)
	}

	fork, err := getFork(user, repoName)
	if err != nil {
		return err
	}

	if fork == nil {
		return simpleClone(options)
	}

	return forkCloner(fork, repoName, options)
}

func simpleClone(options types.CloneOptions) error {
	// git clone  git@github.com:src/repo.git
	output, err := git.Clone(clone.Repository(options.Repo), git.Debug)
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
}

func forkCloner(fork *github.Repository, repoName string, options types.CloneOptions) error {
	// git clone git@github.com:src/repo.git
	output, err := git.Clone(clone.Repository(fork.GetSSHURL()), git.Debug)
	if err != nil {
		log.Println(output)
		return fmt.Errorf("failed to add clone repository %s: %w", fork.GetSSHURL(), err)
	}

	err = os.Chdir(filepath.Clean(repoName))
	if err != nil {
		return err
	}

	// git remote add upstream git@github.com:user/repo.git
	output, err = git.Remote(remote.Add("upstream", options.Repo), git.Debug)
	if err != nil {
		log.Println(output)
		return fmt.Errorf("failed to add remote upstream %s: %w", options.Repo, err)
	}

	// git fetch --multiple origin upstream
	output, err = git.Fetch(fetch.Multiple, fetch.Remote("origin"), fetch.Remote("upstream"), git.Debug)
	if err != nil {
		log.Println(output)
		return fmt.Errorf("failed to add fetch remotes: %w", err)
	}

	return nil
}

func getFork(user, repo string) (*github.Repository, error) {
	ctx := context.Background()
	client := newGitHubClient(ctx)

	ghUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	fork, _, err := client.Repositories.Get(ctx, ghUser.GetLogin(), repo)
	if err != nil {
		v, ok := err.(*github.ErrorResponse)
		if !ok || v != nil && v.Response.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	if fork == nil {
		return createFork(ctx, client, user, repo)
	}

	// fork already exists
	log.Println("fork founded:", fork.GetFullName())

	if !fork.GetFork() {
		return nil, fmt.Errorf("the repository %s is not fork", fork.GetFullName())
	}

	return fork, nil
}

func createFork(ctx context.Context, client *github.Client, user string, repo string) (*github.Repository, error) {
	fmt.Println("No existing fork found.")

	yes, err := choose.Fork()
	if err != nil || !yes {
		return nil, err
	}

	newFork, resp, err := client.Repositories.CreateFork(ctx, user, repo, nil)
	if err != nil {
		_, ok := err.(*github.AcceptedError)
		if !ok || resp == nil || resp.StatusCode != http.StatusAccepted {
			return nil, fmt.Errorf("failed to create a fork: %w", err)
		}
	}

	return newFork, nil
}

func splitUserRepo(u *url.URL) (string, string, error) {
	clean := strings.TrimPrefix(strings.TrimSuffix(u.Path, ".git"), "/")
	parts := strings.Split(clean, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("unable to get user and repositoy name from URL: %v", u)
	}

	return parts[0], parts[1], nil
}
