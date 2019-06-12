package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/v26/github"
	"github.com/ldez/prm/choose"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/local"
	"github.com/ldez/prm/types"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
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
		return 0, errors.Wrap(err, "fail to retrieve pull request from GitHub")
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
		err = pr.Checkout(false)
		if err != nil {
			return err
		}
	} else {
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

		err = config.Save(confs)
		if err != nil {
			return err
		}
	}

	return nil
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

func newGitHubClient(ctx context.Context) *github.Client {
	token := getOrFile("PRM_GITHUB_TOKEN")

	var client *github.Client
	if len(token) == 0 {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}

// getOrFile Attempts to resolve 'key' as an environment variable.
// Failing that, it will check to see if '<key>_FILE' exists.
// If so, it will attempt to read from the referenced file to populate a value.
func getOrFile(envVar string) string {
	envVarValue := os.Getenv(envVar)
	if envVarValue != "" {
		return envVarValue
	}

	fileVar := envVar + "_FILE"
	fileVarValue := os.Getenv(fileVar)
	if fileVarValue == "" {
		return envVarValue
	}

	fileContents, err := ioutil.ReadFile(fileVarValue)
	if err != nil {
		log.Printf("Failed to read the file %s (defined by env var %s): %s", fileVarValue, fileVar, err)
		return ""
	}

	return string(fileContents)
}
