package cmd

import (
	"log"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/remote"
	"github.com/ldez/prm/v3/choose"
	"github.com/ldez/prm/v3/config"
	"github.com/ldez/prm/v3/local"
	"github.com/ldez/prm/v3/types"
)

// InteractiveRemove remove PR.
func InteractiveRemove(conf *config.Configuration) error {
	number, err := choose.PullRequest(conf.PullRequests, true)
	if err != nil || number <= choose.ExitValue {
		return err
	}

	removeOptions := &types.RemoveOptions{}
	if number == choose.AllValue {
		removeOptions.All = true
	} else {
		removeOptions.Numbers = append(removeOptions.Numbers, number)
	}

	return Remove(removeOptions)
}

// Remove remove PR.
func Remove(options *types.RemoveOptions) error {
	// get configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := local.GetGitRepoRoot()
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
		log.Println("remove", pr)

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
			log.Println("remove", pr)

			err := pr.Remove()
			if err != nil {
				return err
			}
		}

		log.Println("remove remote", remoteName)
		out, errRemote := git.Remote(remote.Remove(remoteName), git.Debug)
		if errRemote != nil {
			log.Println(out)
			return errRemote
		}
	}
	return nil
}
