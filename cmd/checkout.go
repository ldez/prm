package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
	"golang.org/x/oauth2"
)

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

	con, err := config.Find(confs, repoDir)
	if err != nil {
		remoteName, err := promptRemoteChoice(remotes)
		if err != nil {
			return err
		}

		confs = append(confs, config.Configuration{
			Directory:  repoDir,
			BaseRemote: remoteName,
		})
		con, err = config.Find(confs, repoDir)
		if err != nil {
			return err
		}
	}

	// check if already exists in config
	for _, pulls := range con.PullRequests {
		for _, pull := range pulls {
			if pull.Number == options.Number {
				return fmt.Errorf("PR already exists: %d", options.Number)
			}
		}
	}

	rmt, err := findRemote(remotes, con.BaseRemote)
	if err != nil {
		return err
	}

	baseRepository, err := newRepository(rmt.URL)
	if err != nil {
		return err
	}

	pr, err := getPullRequest(baseRepository, options.Number)
	if err != nil {
		return err
	}

	// add
	if con.PullRequests == nil {
		con.PullRequests = make(map[string][]types.PullRequest)
	}
	con.PullRequests[pr.Owner] = append(con.PullRequests[pr.Owner], *pr)
	fmt.Println(pr, "checkout")

	err = pr.Checkout()
	if err != nil {
		return err
	}

	err = config.Save(confs)
	if err != nil {
		return err
	}

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

	return remotes
}

func findRemote(remotes []types.Remote, remoteName string) (*types.Remote, error) {
	for _, rmt := range remotes {
		if rmt.Name == remoteName {
			return &rmt, nil
		}
	}
	return nil, fmt.Errorf("Unable to find remote: %s", remoteName)
}

func promptRemoteChoice(remotes []types.Remote) (string, error) {
	for i, rmt := range remotes {
		fmt.Printf("%d: %s (%s)\n", i, rmt.Name, rmt.URL)
	}
	fmt.Println("Choose the base remote:")

	reader := bufio.NewReader(os.Stdin)
	rawAnswer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer, err := strconv.ParseInt(strings.TrimSpace(rawAnswer), 10, 8)
	if err != nil || answer > int64(len(remotes)) || answer < 0 {
		return "", fmt.Errorf("Invalid answer: %s", rawAnswer)
	}

	return remotes[answer].Name, nil
}

func newRepository(URL string) (*types.Repository, error) {
	exp := regexp.MustCompile(`(?:git@github.com:|https://github.com/)([^/]+)/(.+).git`)

	parts := exp.FindStringSubmatch(URL)

	if len(parts) < 3 {
		return nil, fmt.Errorf("Invalid URL: %s", URL)
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
		Owner:      *pr.Head.Repo.Owner.Login,
		BranchName: *pr.Head.Ref,
		Number:     number,
	}, nil
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
