package config

import (
	"os/user"
	"path/filepath"
)

const defaultFileName = ".prm"

// GetPath get the configuration file path.
func GetPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, defaultFileName), nil
}
