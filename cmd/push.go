package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/revparse"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

// Push push to the PR branch.
func Push(options *types.PushOptions) error {

	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}

	con, err := config.Find(confs, repoDir)
	if err != nil {
		return err
	}

	number, err := getPRNUmber(options.Number)
	if err != nil {
		return err
	}

	pr, err := con.FindPullRequests(number)
	if err != nil {
		return err
	}

	fmt.Println("push", pr)

	err = pr.Push(options.Force)
	if err != nil {
		return err
	}

	return nil
}

func getPRNUmber(manualNumber int) (int, error) {
	if manualNumber == 0 {
		out, err := git.RevParse(revparse.AbbrevRef(""), revparse.Args("HEAD"))
		if err != nil {
			log.Println(out)
			return 0, err
		}

		number, err := parsePRNumber(out)
		if err != nil {
			return 0, err
		}
		return number, nil
	} else {
		return manualNumber, nil
	}
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

	return 0, fmt.Errorf("Unable to parse: %s", out)
}
