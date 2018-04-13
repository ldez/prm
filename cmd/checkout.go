package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/choose"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Checkout checkout a PR.
// TODO simplify this function
// nolint: gocyclo
func Checkout(options *types.CheckoutOptions) error {
	// Get all remotes
	output, err := git.Remote(remote.Verbose, git.Debug)
	if err != nil {
		log.Println(output)
		return err
	}

	remotes := getRemotes(output)

	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}

	conf, err := config.Find(confs, repoDir)
	if err != nil {
		var remoteName string
		if len(remotes) == 1 {
			remoteName = remotes[0].Name
		} else {
			remoteName, err = promptRemoteChoice(remotes)
			if err != nil {
				return err
			}
		}

		confs = append(confs, config.Configuration{
			Directory:  repoDir,
			BaseRemote: remoteName,
		})
		conf, err = config.Find(confs, repoDir)
		if err != nil {
			return err
		}
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
		// remote checkout
		rmt, err := findRemote(remotes, conf.BaseRemote)
		if err != nil {
			return err
		}

		baseRepository, err := newRepository(rmt.URL)
		if err != nil {
			return err
		}

		prNumber, err := getPullRequestNumber(options, baseRepository)
		if err != nil || prNumber <= 0 || prNumber == choose.ExitValue {
			return err
		}

		pr, err = getPullRequest(baseRepository, prNumber)
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

	fmt.Println("checkout", pr)

	return nil
}

func getRemotes(output string) []types.Remote {
	lines := strings.Split(output, "\n")

	remoteMap := make(map[string]types.Remote)

	for _, line := range lines {
		if len(line) != 0 {
			exp := regexp.MustCompile(`[\t|\s{2,}]`)
			elt := exp.Split(line, 3)

			name := elt[0]
			rmt := types.Remote{
				Name: name,
				URL:  elt[1],
			}
			remoteMap[name] = rmt
		}
	}

	var remotes []types.Remote
	for _, entry := range remoteMap {
		remotes = append(remotes, entry)
	}

	sort.Sort(types.ByRemoteName(remotes))

	return remotes
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

func findRemote(remotes []types.Remote, remoteName string) (*types.Remote, error) {
	for _, rmt := range remotes {
		if rmt.Name == remoteName {
			return &rmt, nil
		}
	}
	return nil, fmt.Errorf("unable to find remote: %s", remoteName)
}

// TODO use survey
func promptRemoteChoice(remotes []types.Remote) (string, error) {
	for i, rmt := range remotes {
		fmt.Printf("%d: %s (%s)\n", i, rmt.Name, rmt.URL)
	}
	fmt.Println("Choose the remote related to PR (main remote):")

	reader := bufio.NewReader(os.Stdin)
	rawAnswer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer, err := strconv.ParseInt(strings.TrimSpace(rawAnswer), 10, 8)
	if err != nil || answer > int64(len(remotes)) || answer < 0 {
		return "", fmt.Errorf("invalid answer: %s", rawAnswer)
	}

	return remotes[answer].Name, nil
}

func newRepository(URL string) (*types.Repository, error) {
	exp := regexp.MustCompile(`(?:git@github.com:|https://github.com/)([^/]+)/(.+)\.git`)

	if !exp.MatchString(URL) {
		return nil, fmt.Errorf("invalid URL: %s", URL)
	}

	parts := exp.FindStringSubmatch(URL)

	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URL: %s", URL)
	}

	return &types.Repository{
		Owner: parts[1],
		Name:  parts[2],
	}, nil
}

func getPullRequest(baseRepository *types.Repository, number int) (*types.PullRequest, error) {
	ctx := context.Background()
	client := newGitHubClient(ctx, "")

	pr, _, err := client.PullRequests.Get(ctx, baseRepository.Owner, baseRepository.Name, number)
	if err != nil {
		return nil, err
	}

	return &types.PullRequest{
		Project:    baseRepository.Name,
		Owner:      pr.Head.Repo.Owner.GetLogin(),
		BranchName: pr.Head.GetRef(),
		Number:     number,
	}, nil
}

func getPullRequestNumber(options *types.CheckoutOptions, baseRepository *types.Repository) (int, error) {
	prNumber := options.Number
	if prNumber == 0 {
		ctx := context.Background()
		client := newGitHubClient(ctx, "")

		opt := &github.PullRequestListOptions{
			State:       "open",
			ListOptions: github.ListOptions{PerPage: 25},
		}

		prs, _, err := client.PullRequests.List(ctx, baseRepository.Owner, baseRepository.Name, opt)
		if err != nil {
			return 0, errors.Wrap(err, "fail to retrieve pull request from GitHub")
		}

		return choose.RemotePulRequest(prs)
	}

	return prNumber, nil
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
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
