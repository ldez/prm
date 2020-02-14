package types

import (
	"fmt"
	"strings"

	"github.com/ldez/prm/local"
	giturls "github.com/whilp/git-urls"
)

// Repository Git repository model.
type Repository struct {
	Owner string
	Name  string
}

func newRepository(uri string) (*Repository, error) {
	u, err := giturls.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s: %w", uri, err)
	}

	parts := strings.Split(strings.TrimSuffix(strings.TrimSuffix(u.Path, "/"), ".git"), "/")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid URL: %s", uri)
	}

	return &Repository{Owner: parts[len(parts)-2], Name: parts[len(parts)-1]}, nil
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
