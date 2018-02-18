package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/types"
)

// Remove remove PR.
func Remove(options *types.RemoveOptions) error {
	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return err
	}

	conf, err := config.Find(confs, repoDir)
	if err != nil {
		return err
	}

	if options.All {
		err = removeAll(conf)
		if err != nil {
			return err
		}
		conf.PullRequests = make(map[string][]types.PullRequest)
	} else {
		for _, prNumber := range options.Numbers {
			err = removePR(conf, prNumber)
			if err != nil {
				return err
			}
		}
	}

	return config.Save(confs)
}

func removePR(conf *config.Configuration, prNumber int) error {
	if pr, err := conf.FindPullRequests(prNumber); err == nil {
		fmt.Println("remove", pr)

		err = pr.Remove()
		if err != nil {
			return err
		}

		if conf.RemovePullRequest(pr) == 0 {
			err = pr.RemoveRemote()
			if err != nil {
				return err
			}
		}
	} else {
		log.Println(err)
	}
	return nil
}

func removeAll(conf *config.Configuration) error {
	for remoteName, prs := range conf.PullRequests {
		for _, pr := range prs {
			fmt.Println("remove", pr)

			err := pr.Remove()
			if err != nil {
				return err
			}
		}

		fmt.Println("remove remote", remoteName)
		out, errRemote := git.Remote(remote.Remove(remoteName), git.Debug)
		if errRemote != nil {
			log.Println(out)
			return errRemote
		}
	}
	return nil
}
