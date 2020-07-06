package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v32/github"
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
	user, repoName, err := splitUserRepo(options.Repo)
	if err != nil {
		return err
	}

	if options.UserAsRootDir {
		if err = os.MkdirAll(filepath.Clean(user), 0o755); err != nil {
			return err
		}

		if err = os.Chdir(filepath.Clean(user)); err != nil {
			return err
		}
	}

	if options.NoFork {
		return cloneSimple(options)
	}

	ctx := context.Background()
	cl := newCloner(ctx)

	forkUsername, err := cl.getForkUser(ctx)
	if err != nil {
		return err
	}

	if forkUsername == "" {
		return cloneSimple(options)
	}

	fork, err := cl.searchFork(ctx, forkUsername, user, repoName)
	if err != nil {
		return err
	}

	// fork already exists
	if fork != nil {
		log.Println("fork founded:", fork.GetFullName())
		return cloneFork(fork, repoName, options)
	}

	fork, err = cl.createFork(ctx, user, repoName, options.Organization)
	if err != nil {
		return err
	}

	if fork == nil {
		return cloneSimple(options)
	}

	return cloneFork(fork, repoName, options)
}

type cloner struct {
	client *github.Client
}

func newCloner(ctx context.Context) cloner {
	return cloner{
		client: newGitHubClient(ctx),
	}
}

func (c cloner) getForkUser(ctx context.Context) (string, error) {
	if hasToken() {
		authUser, _, err := c.client.Users.Get(ctx, "")
		if err != nil {
			return "", err
		}

		return authUser.GetLogin(), nil
	}

	fmt.Println("---------------------------------------------------------")
	fmt.Printf("Set %s or %s to allow PRM to detect automatically your username:\n", tokenEnvVar, tokenEnvVar+fileSuffixEnvVar)
	fmt.Println("- https://ldez.github.io/prm/#prm-github-token")
	fmt.Println("- https://ldez.github.io/prm/#prm-github-token-file")
	fmt.Println("---------------------------------------------------------")

	return choose.Username()
}

func (c cloner) searchFork(ctx context.Context, me, user, repoName string) (*github.Repository, error) {
	query := fmt.Sprintf("user:%s fork:only in:name %s", me, repoName)

	searchResult, _, err := c.client.Search.Repositories(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	if searchResult.GetTotal() == 0 {
		return nil, nil
	}

	for _, repository := range searchResult.Repositories {
		if repository.GetFork() {
			repo, _, err := c.client.Repositories.Get(ctx, me, repository.GetName())
			if err != nil {
				return nil, err
			}

			srcRepoFullName := fmt.Sprintf("%s/%s", user, repoName)

			if repo.GetParent().GetFullName() == srcRepoFullName {
				return repo, nil
			}

			if repo.GetSource().GetFullName() == srcRepoFullName {
				return repo, nil
			}
		}
	}

	return nil, nil
}

func (c cloner) createFork(ctx context.Context, user, repo, org string) (*github.Repository, error) {
	if !hasToken() {
		fmt.Println("---------------------------------------------------------")
		fmt.Printf("Set %s or %s to allow to fork automatically:\n", tokenEnvVar, tokenEnvVar+fileSuffixEnvVar)
		fmt.Println("- https://ldez.github.io/prm/#prm-github-token")
		fmt.Println("- https://ldez.github.io/prm/#prm-github-token-file")
		fmt.Println("---------------------------------------------------------")
		return nil, nil
	}

	fmt.Println("No existing fork found: creating a new fork.")

	yes, err := choose.Fork()
	if err != nil || !yes {
		return nil, err
	}

	var opt *github.RepositoryCreateForkOptions
	if org == "" {
		opt = &github.RepositoryCreateForkOptions{Organization: org}
	}

	newFork, resp, err := c.client.Repositories.CreateFork(ctx, user, repo, opt)
	if err != nil {
		_, ok := err.(*github.AcceptedError)
		if !ok || resp == nil || resp.StatusCode != http.StatusAccepted {
			return nil, fmt.Errorf("failed to create a fork: %w", err)
		}
	}

	return newFork, nil
}

func cloneSimple(options types.CloneOptions) error {
	// git clone  git@github.com:src/repo.git
	output, err := git.Clone(clone.Repository(options.Repo), git.Debug)
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
}

func cloneFork(fork *github.Repository, repoName string, options types.CloneOptions) error {
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

func splitUserRepo(rawURL string) (string, string, error) {
	u, err := giturls.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("unable to get user and repositoy name from URL (%s): %w", rawURL, err)
	}

	clean := strings.TrimPrefix(strings.TrimSuffix(u.Path, ".git"), "/")
	parts := strings.Split(clean, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("unable to get user and repositoy name from URL: %s", rawURL)
	}

	return parts[0], parts[1], nil
}
