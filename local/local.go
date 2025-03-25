package local

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/ldez/go-git-cmd-wrapper/v2/branch"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/remote"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
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

// ByRemoteName sorts remote by name.
type ByRemoteName Remotes

func (r ByRemoteName) Len() int           { return len(r) }
func (r ByRemoteName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRemoteName) Less(i, j int) bool { return r[i].Name < r[j].Name }

// GetCurrentPRNumber gets the current PR number.
func GetCurrentPRNumber(manualNumber int) (int, error) {
	if manualNumber == 0 {
		return GetCurrentBranchPRNumber()
	}

	return manualNumber, nil
}

// GetCurrentBranchPRNumber gets the current branch PR number.
func GetCurrentBranchPRNumber() (int, error) {
	output, err := GetCurrentBranchName()
	if err != nil {
		log.Print(output)
		return 0, err
	}

	return parsePRNumber(output)
}

// GetCurrentBranchName gets the current branch name.
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

// GetGitRepoRoot gets the root of the git repository.
func GetGitRepoRoot() (string, error) {
	output, err := git.RevParse(revparse.ShowToplevel)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// GetRemotes gets git remotes.
func GetRemotes() (Remotes, error) {
	output, err := git.Remote(remote.Verbose, git.Debug)
	if err != nil {
		log.Print(output)
		return nil, err
	}

	return parseRemotes(output), nil
}

// GetBranches gets git branches.
func GetBranches() ([]string, error) {
	output, err := git.Branch(branch.Format("%(refname:short)"), branch.Sort("refname"), git.Debug)
	if err != nil {
		log.Print(output)
		return nil, err
	}

	return parseBranches(output), nil
}

func parseBranches(output string) []string {
	var branches []string

	for _, name := range strings.Split(output, "\n") {
		b := strings.TrimSpace(name)
		if b == "" {
			continue
		}

		branches = append(branches, b)
	}

	sort.Slice(branches, func(i, _ int) bool {
		switch branches[i] {
		case "main", "master":
			return true

		default:
			return false
		}
	})

	return branches
}

func parseRemotes(output string) Remotes {
	lines := strings.Split(output, "\n")

	remoteMap := make(map[string]Remote)

	for _, line := range lines {
		if line != "" {
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
