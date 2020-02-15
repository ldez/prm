package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

const (
	// GitHub token.
	tokenEnvVar = "PRM_GITHUB_TOKEN"
	// GitHub Enterprise API base URL.
	apiBaseURLEnvVar = "PRM_GITHUB_API_BASE_URL"
)

// HasToken checks if the GitHub token is present.
func HasToken() bool {
	return getOrFile(tokenEnvVar) != ""
}

func newGitHubClient(ctx context.Context) *github.Client {
	token := getOrFile(tokenEnvVar)

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

	baseURL := getOrFile(apiBaseURLEnvVar)
	if baseURL != "" {
		var err error
		client.BaseURL, err = url.Parse(strings.TrimSuffix(baseURL, "/") + "/")
		if err != nil {
			panic(fmt.Sprintf("invalid domain endpoint: %v", err))
		}
	}

	return client
}
