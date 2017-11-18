package local

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/revparse"
)

// GetCurrentPRNumber get the current PR number.
func GetCurrentPRNumber(manualNumber int) (int, error) {
	if manualNumber == 0 {
		return GetCurrentBranchPRNumber()
	}
	return manualNumber, nil
}

// GetCurrentBranchPRNumber get the current branch PR number.
func GetCurrentBranchPRNumber() (int, error) {
	out, err := GetCurrentBranchName()
	if err != nil {
		log.Println(out)
		return 0, err
	}

	return parsePRNumber(out)
}

// GetCurrentBranchName get the current branch name.
func GetCurrentBranchName() (string, error) {
	out, err := git.RevParse(revparse.AbbrevRef(""), revparse.Args("HEAD"))
	if err != nil {
		log.Println(out)
		return "", err
	}
	return strings.TrimSpace(out), nil
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
