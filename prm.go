package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/containous/flaeg"
	"github.com/google/go-github/github"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
	"golang.org/x/oauth2"
)

func main() {

	emptyConfig := &types.NoOption{}
	rootCmd := &flaeg.Command{
		Name:                  "prm",
		Description:           "PRM - The Pull Request Manager.",
		Config:                emptyConfig,
		DefaultPointersConfig: &types.NoOption{},
		Run: func() error {
			return nil
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])

	// Checkout

	checkoutOptions := &types.CheckoutOptions{}

	checkoutCmd := &flaeg.Command{
		Name:                  "c",
		Description:           "Checkout a PR.",
		Config:                checkoutOptions,
		DefaultPointersConfig: &types.CheckoutOptions{},
	}
	checkoutCmd.Run = func() error {
		err := requirePRNumber(checkoutOptions.Number, checkoutCmd.Name)
		if err != nil {
			log.Fatalln(err)
		}
		err = checkoutCommand(checkoutOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	flag.AddCommand(checkoutCmd)

	// Remove

	removeOptions := &types.RemoveOptions{}

	removeCmd := &flaeg.Command{
		Name:                  "rm",
		Description:           "Remove one or more PRs from the local repository.",
		Config:                removeOptions,
		DefaultPointersConfig: &types.RemoveOptions{},
	}
	removeCmd.Run = func() error {
		err := requirePRNumber(removeOptions.Number, removeCmd.Name)
		if !removeOptions.All && err != nil {
			log.Fatalln(err)
		}
		err = removeCommand(removeOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	flag.AddCommand(removeCmd)

	// Push Force

	pushForceOptions := &types.PushForceOptions{}

	pushForceCmd := &flaeg.Command{
		Name:                  "pf",
		Description:           "Push force a PR.",
		Config:                pushForceOptions,
		DefaultPointersConfig: &types.PushForceOptions{},
	}
	pushForceCmd.Run = func() error {
		err := requirePRNumber(pushForceOptions.Number, pushForceCmd.Name)
		if err != nil {
			log.Fatalln(err)
		}
		err = pushForceCommand(pushForceOptions)
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	flag.AddCommand(pushForceCmd)

	// List

	listOptions := &types.ListOptions{}

	listCmd := &flaeg.Command{
		Name:                  "list",
		Description:           "Display all current PRs.",
		Config:                listOptions,
		DefaultPointersConfig: &types.ListOptions{},
		Run: func() error {
			err := listCommand(listOptions)
			if err != nil {
				log.Println(err)
			}
			return nil
		},
	}

	flag.AddCommand(listCmd)

	// Run command
	flag.Run()
}

func listCommand(options *types.ListOptions) error {
	// get configuration
	configs, err := config.ReadFile()
	if err != nil {
		return err
	}

	if options.All {
		for _, conf := range configs {
			fmt.Println(conf.Directory)
			for _, prs := range conf.PullRequests {
				for _, pr := range prs {
					fmt.Printf("* %d: %s - %s\n", pr.Number, pr.Owner, pr.BranchName)
				}
			}
		}
	} else {
		repoDir, err := os.Getwd()
		if err != nil {
			return err
		}

		conf, err := config.Find(configs, repoDir)
		if err != nil {
			return err
		}

		for _, prs := range conf.PullRequests {
			for _, pr := range prs {
				fmt.Printf("* %d: %s - %s\n", pr.Number, pr.Owner, pr.BranchName)
			}
		}
	}

	return nil
}

func checkoutCommand(options *types.CheckoutOptions) error {

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

	//pr.Checkout()

	err = config.Save(confs)
	if err != nil {
		return err
	}

	return nil
}

func removeCommand(options *types.RemoveOptions) error {

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
		return err
	}

	if options.All {
		for _, prs := range con.PullRequests {
			for _, pr := range prs {
				fmt.Println(pr, "remove")
				//pr.RemoveRemote()
				//pr.Remove()
			}
		}
		con.PullRequests = make(map[string][]types.PullRequest)
	} else {
		if pr, err := con.FindPullRequests(options.Number); err == nil {
			fmt.Println(pr, "remove")
			// pr.Remove()
			if con.RemovePullRequest(pr) == 0 {
				// pr.RemoveRemote()
			}
		} else {
			log.Println(err)
		}
	}

	err = config.Save(confs)
	if err != nil {
		return err
	}

	return nil
}

func pushForceCommand(options *types.PushForceOptions) error {

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
		return err
	}

	if pr, err := con.FindPullRequests(options.Number); err == nil {
		fmt.Println(pr, "push force")
		//pr.PushForce()
	}

	return nil
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

func newRepository(URL string) (*types.Repository, error) {
	// https://github.com/ldez/prm.git
	// git@github.com:containous/traefik.git
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
		Owner:      *pr.Head.Repo.Owner.Login,
		BranchName: *pr.Head.Ref,
		Number:     number,
	}, nil
}

func requirePRNumber(number int, action string) error {
	if number <= 0 {
		return fmt.Errorf("You must provide a PR number. ex: 'prm %s --number=1235'", action)
	}
	return nil
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
