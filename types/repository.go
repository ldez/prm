package types

import (
	"fmt"
	"regexp"

	"github.com/ldez/prm/local"
)

// Repository Git repository model.
type Repository struct {
	Owner string
	Name  string
}

func newRepository(URL string) (*Repository, error) {
	exp := regexp.MustCompile(`(?:git@github.com:|https://github.com/)([^/]+)/(.+)\.git`)

	if !exp.MatchString(URL) {
		return nil, fmt.Errorf("invalid URL: %s", URL)
	}

	parts := exp.FindStringSubmatch(URL)

	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URL: %s", URL)
	}

	return &Repository{
		Owner: parts[1],
		Name:  parts[2],
	}, nil
}

// GetRepository get repository information by remote name.
func GetRepository(baseRemote string) (*Repository, error) {
	remotes, err := local.GetRemotes()
	if err != nil {
		return nil, err
	}

	// remote checkout
	rmt, err := remotes.Find(baseRemote)
	if err != nil {
		return nil, err
	}

	return newRepository(rmt.URL)
}
