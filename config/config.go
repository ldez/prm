package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/ldez/prm/types"
)

// Configuration is the global application configuration model.
type Configuration struct {
	Directory    string                         `json:"directory,omitempty"`
	BaseRemote   string                         `json:"base_remote,omitempty"`
	PullRequests map[string][]types.PullRequest `json:"pull_requests,omitempty"`
}

const defaultFileName = ".prm"

var getPathFunc = GetPath

// RemovePullRequest remove a pull request.
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

// FindPullRequests find a pull request by number.
func (c *Configuration) FindPullRequests(number int) (*types.PullRequest, error) {
	for _, prs := range c.PullRequests {
		for _, pr := range prs {
			if pr.Number == number {
				return &pr, nil
			}
		}
	}

	return nil, fmt.Errorf("Unable to find PR: %d", number)
}

// Find find a configuration by directory name.
func Find(configurations []Configuration, directory string) (*Configuration, error) {

	for i, config := range configurations {
		if config.Directory == directory {
			return &configurations[i], nil
		}
	}

	return nil, fmt.Errorf("No existing configuration for: %s", directory)
}

// ReadFile read the configuration file and load the configuration into an array.
func ReadFile() ([]Configuration, error) {

	var configs []Configuration

	filePath, err := getPathFunc()
	if err != nil {
		return configs, err
	}

	if _, err := os.Stat(filePath); err != nil {
		log.New(os.Stdout, "INFO: ", log.LstdFlags).Printf("Create the configuration file: %s", filePath)

		content, err := json.MarshalIndent(configs, "", "  ")
		if err != nil {
			return configs, err
		}

		file, err := os.Create(filePath)
		if err != nil {
			return configs, err
		}

		_, err = file.Write(content)
		if err != nil {
			return configs, err
		}

		defer file.Close()
	}

	file, err := ioutil.ReadFile(filePath)

	err = json.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// Save save the configuration into a file.
func Save(configs []Configuration) error {

	filePath, err := getPathFunc()
	if err != nil {
		return err
	}

	confJSON, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, confJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetPath get the configuration file path.
func GetPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, defaultFileName), nil
}
