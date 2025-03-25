package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ldez/prm/v3/local"
	"github.com/ldez/prm/v3/types"
)

// Configuration is the global application configuration model.
type Configuration struct {
	PullRequests map[string][]types.PullRequest `json:"pull_requests,omitempty"`
	Directory    string                         `json:"directory,omitempty"`
	BaseRemote   string                         `json:"base_remote,omitempty"`
	MainBranch   string                         `json:"main_branch,omitempty"`
}

var getPathFunc = GetPath

// RemovePullRequest removes a pull request.
func (c *Configuration) RemovePullRequest(pull *types.PullRequest) int {
	prs := c.PullRequests[pull.Owner]

	index := c.findPullRequestIndex(pull)

	var b []types.PullRequest

	if index != -1 {
		b = append(b, prs[:index]...)
		b = append(b, prs[index+1:]...)

		if len(b) == 0 {
			// It's the only PR for this owner
			delete(c.PullRequests, pull.Owner)
		} else {
			c.PullRequests[pull.Owner] = b
		}
	}

	return len(b)
}

func (c *Configuration) findPullRequestIndex(pull *types.PullRequest) int {
	prs := c.PullRequests[pull.Owner]

	for i, pr := range prs {
		if pr.Number == pull.Number {
			return i
		}
	}

	return -1
}

// FindPullRequests finds a pull request by number.
func (c *Configuration) FindPullRequests(number int) (*types.PullRequest, error) {
	for _, prs := range c.PullRequests {
		for _, pr := range prs {
			if pr.Number == number {
				return &pr, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find PR: %d", number)
}

// Get configuration for the current directory.
func Get() (*Configuration, error) {
	confs, err := ReadFile()
	if err != nil {
		return nil, err
	}

	repoDir, err := local.GetGitRepoRoot()
	if err != nil {
		return nil, err
	}

	return Find(confs, repoDir)
}

// Find finds a configuration by directory name.
func Find(configurations []Configuration, directory string) (*Configuration, error) {
	for i, config := range configurations {
		if config.Directory == directory {
			return &configurations[i], nil
		}
	}

	return nil, fmt.Errorf("no existing configuration for: %s", directory)
}

// ReadFile reads or creates the configuration file and load the configuration into an array.
func ReadFile() ([]Configuration, error) {
	var configs []Configuration

	filePath, err := getPathFunc()
	if err != nil {
		return configs, err
	}

	if _, errStat := os.Stat(filePath); errStat != nil {
		log.New(os.Stdout, "INFO: ", log.LstdFlags).Printf("Create the configuration file: %s", filePath)

		content, errMarshal := json.MarshalIndent(configs, "", "  ")
		if errMarshal != nil {
			return configs, errMarshal
		}

		errDir := createDirectory(filePath)
		if errDir != nil {
			return configs, errDir
		}

		file, errCreate := os.Create(filePath)
		if errCreate != nil {
			return configs, errCreate
		}

		defer func() {
			errClose := file.Close()
			if errClose != nil {
				log.Println(errClose)
			}
		}()

		_, errWrite := file.Write(content)
		if errWrite != nil {
			return configs, errWrite
		}
	}

	file, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// Save saves the configuration into a file.
func Save(configs []Configuration) error {
	filePath, err := getPathFunc()
	if err != nil {
		return err
	}

	slices.SortFunc(configs, func(a, b Configuration) int {
		return strings.Compare(a.Directory, b.Directory)
	})

	confJSON, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	err = createDirectory(filePath)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, confJSON, 0o644)
}

func createDirectory(filePath string) error {
	baseDir := path.Dir(filePath)

	if _, errDirStat := os.Stat(baseDir); errDirStat != nil {
		errDir := os.MkdirAll(baseDir, 0o700)
		if errDir != nil {
			return errDir
		}
	}

	return nil
}
