package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/revparse"
)

func getPRNumber(manualNumber int) (int, error) {
	if manualNumber == 0 {
		return getBranchPRNumber()
	}
	return manualNumber, nil
}

func getBranchPRNumber() (int, error) {
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
