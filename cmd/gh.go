package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

const (
	// GitHub token.
	tokenEnvVar = "PRM_GITHUB_TOKEN"
	// File suffix.
	fileSuffixEnvVar = "_FILE"
	// GitHub Enterprise API base URL.
	apiBaseURLEnvVar = "PRM_GITHUB_API_BASE_URL"
)

// hasToken checks if the GitHub token is present.
func hasToken() bool {
	return getOrFile(tokenEnvVar) != ""
}

func newGitHubClient(ctx context.Context) *github.Client {
	token := getOrFile(tokenEnvVar)

	var client *github.Client
	if token == "" {
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

// getOrFile Attempts to resolve 'key' as an environment variable.
// Failing that, it will check to see if '<key>_FILE' exists.
// If so, it will attempt to read from the referenced file to populate a value.
func getOrFile(envVar string) string {
	envVarValue := os.Getenv(envVar)
	if envVarValue != "" {
		return envVarValue
	}

	fileVar := envVar + fileSuffixEnvVar

	fileVarValue := os.Getenv(fileVar)
	if fileVarValue == "" {
		return envVarValue
	}

	fileContents, err := os.ReadFile(filepath.Clean(fileVarValue))
	if err != nil {
		log.Printf("Failed to read the file %q (defined by env var %q): %v", fileVarValue, fileVar, err)
		return ""
	}

	return string(bytes.TrimSpace(fileContents))
}
