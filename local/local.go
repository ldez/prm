package local

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/go-git-cmd-wrapper/revparse"
)

// Remote Git remote model.
type Remote struct {
	Name string
	URL  string
}

// Remotes a list of Remote.
type Remotes []Remote

// Find a remote by name.
func (r Remotes) Find(remoteName string) (*Remote, error) {
	for _, rmt := range r {
		if rmt.Name == remoteName {
			return &rmt, nil
		}
	}

	return nil, fmt.Errorf("unable to find remote %q in %v", remoteName, r)
}

// ByRemoteName sort remote by name.
type ByRemoteName Remotes

func (r ByRemoteName) Len() int           { return len(r) }
func (r ByRemoteName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRemoteName) Less(i, j int) bool { return r[i].Name < r[j].Name }

// GetCurrentPRNumber get the current PR number.
func GetCurrentPRNumber(manualNumber int) (int, error) {
	if manualNumber == 0 {
		return GetCurrentBranchPRNumber()
	}
	return manualNumber, nil
}

// GetCurrentBranchPRNumber get the current branch PR number.
func GetCurrentBranchPRNumber() (int, error) {
	output, err := GetCurrentBranchName()
	if err != nil {
		log.Print(output)
		return 0, err
	}

	return parsePRNumber(output)
}

// GetCurrentBranchName get the current branch name.
func GetCurrentBranchName() (string, error) {
	output, err := git.RevParse(revparse.AbbrevRef(""), revparse.Args("HEAD"))
	if err != nil {
		log.Print(output)
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func parsePRNumber(out string) (int, error) {
	exp := regexp.MustCompile(`(\d+)--.+`)
	parts := exp.FindStringSubmatch(out)

	if len(parts) == 2 {
		number, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return 0, err
		}

		return int(number), nil
	}

	return 0, fmt.Errorf("unable to parse: %s", out)
}

// GetGitRepoRoot get the root of the git repository.
func GetGitRepoRoot() (string, error) {
	output, err := git.RevParse(revparse.ShowToplevel)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// GetRemotes get git remotes.
func GetRemotes() (Remotes, error) {
	output, err := git.Remote(remote.Verbose, git.Debug)
	if err != nil {
		log.Print(output)
		return nil, err
	}

	return parseRemotes(output), nil
}

func parseRemotes(output string) Remotes {
	lines := strings.Split(output, "\n")

	remoteMap := make(map[string]Remote)

	for _, line := range lines {
		if len(line) != 0 {
			elt := strings.FieldsFunc(line, unicode.IsSpace)

			name := elt[0]
			rmt := Remote{
				Name: name,
				URL:  elt[1],
			}
			remoteMap[name] = rmt
		}
	}

	var remotes Remotes
	for _, entry := range remoteMap {
		remotes = append(remotes, entry)
	}

	sort.Sort(ByRemoteName(remotes))

	return remotes
}
