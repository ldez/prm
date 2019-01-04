package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ldez/prm/local"
)

// Repository Git repository model.
type Repository struct {
	Owner string
	Name  string
}

func newRepository(URL string) (*Repository, error) {
	exp := regexp.MustCompile(`(?:git@github.com:|https://github.com/)([^/]+)/(.+)`)

	if !exp.MatchString(URL) {
		return nil, fmt.Errorf("invalid URL: %s", URL)
	}

	parts := exp.FindStringSubmatch(URL)

	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid URL: %s", URL)
	}

	name := strings.TrimSuffix(strings.TrimSuffix(parts[2], ".git"), "/")
	return &Repository{Owner: parts[1], Name: name}, nil
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
