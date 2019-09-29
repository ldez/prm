package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/containous/flaeg"
	"github.com/ldez/prm/choose"
	"github.com/ldez/prm/cmd"
	"github.com/ldez/prm/config"
	"github.com/ldez/prm/local"
	"github.com/ldez/prm/meta"
	"github.com/ldez/prm/types"
	"github.com/ogier/pflag"
	"github.com/pkg/errors"
)

func main() {
	rootCmd := &flaeg.Command{
		Name:                  "prm",
		Description:           "PRM - The Pull Request Manager.",
		Run:                   safe(rootRun),
		Config:                &types.NoOption{},
		DefaultPointersConfig: &types.NoOption{},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])

	flag.AddParser(reflect.TypeOf(types.PRNumbers{}), &types.PRNumbers{})

	// Checkout
	flag.AddCommand(createCheckout())

	// Remove
	flag.AddCommand(createRemove())

	// Push Force
	flag.AddCommand(createPushForce())

	// Push
	flag.AddCommand(createPush())

	// Pull
	flag.AddCommand(createPull())

	// List
	flag.AddCommand(createList())

	// version
	flag.AddCommand(createVersion())

	// Run command
	err := flag.Run()
	if err != nil && err != pflag.ErrHelp {
		log.Printf("Error: %v\n", err)
	}
}

func rootRun() error {
	conf, err := config.Get()
	if err != nil {
		return err
	}

	action, err := choose.Action()
	if err != nil {
		return err
	}

	switch action {
	case choose.ActionList:
		return cmd.Switch(&types.ListOptions{})
	case choose.ActionCheckout:
		return cmd.InteractiveCheckout(conf)
	case choose.ActionRemove:
		return cmd.InteractiveRemove(conf)
	case choose.ActionProjects:
		return cmd.Switch(&types.ListOptions{All: true})
	}

	return nil
}

func createCheckout() *flaeg.Command {
	checkoutOptions := &types.CheckoutOptions{}

	checkoutCmd := &flaeg.Command{
		Name:                  "c",
		Description:           "Checkout a PR (create a local branch and add remote).",
		Config:                checkoutOptions,
		DefaultPointersConfig: &types.CheckoutOptions{},
	}
	checkoutCmd.Run = safe(func() error {
		if checkoutOptions.Number != 0 {
			return cmd.Checkout(checkoutOptions)
		}

		conf, err := config.Get()
		if err != nil {
			return err
		}
		return cmd.InteractiveCheckout(conf)
	})

	return checkoutCmd
}

func createRemove() *flaeg.Command {
	removeOptions := &types.RemoveOptions{}

	removeCmd := &flaeg.Command{
		Name:                  "rm",
		Description:           "Remove one or more PRs from the current local repository.",
		Config:                removeOptions,
		DefaultPointersConfig: &types.RemoveOptions{},
	}
	removeCmd.Run = safe(removeRun(removeCmd.Name, removeOptions))

	return removeCmd
}

func removeRun(action string, removeOptions *types.RemoveOptions) func() error {
	return func() error {
		if removeOptions.All {
			return cmd.Remove(removeOptions)
		}

		if !removeOptions.NoPrompt && len(removeOptions.Numbers) == 0 {
			conf, err := config.Get()
			if err != nil {
				return err
			}

			return cmd.InteractiveRemove(conf)
		}

		err := requirePRNumbers(removeOptions.Numbers, action)
		if err != nil {
			return err
		}

		return cmd.Remove(removeOptions)
	}
}

func createPushForce() *flaeg.Command {
	pushForceOptions := &types.PushOptions{Force: true}

	pushForceCmd := &flaeg.Command{
		Name:                  "pf",
		Description:           "Push force to the PR branch.",
		Config:                pushForceOptions,
		DefaultPointersConfig: &types.PushOptions{},
	}
	pushForceCmd.Run = safe(func() error {
		return cmd.Push(pushForceOptions)
	})

	return pushForceCmd
}

func createPush() *flaeg.Command {
	pushOptions := &types.PushOptions{}

	pushCmd := &flaeg.Command{
		Name:                  "push",
		Description:           "Push to the PR branch.",
		Config:                pushOptions,
		DefaultPointersConfig: &types.PushOptions{},
	}
	pushCmd.Run = safe(func() error {
		return cmd.Push(pushOptions)
	})

	return pushCmd
}

func createPull() *flaeg.Command {
	pullOptions := &types.PullOptions{}
	pullCmd := &flaeg.Command{
		Name:                  "pull",
		Description:           "Pull to the PR branch.",
		Config:                pullOptions,
		DefaultPointersConfig: &types.PullOptions{},
	}
	pullCmd.Run = safe(func() error {
		return cmd.Pull(pullOptions)
	})
	return pullCmd
}

func createList() *flaeg.Command {
	listOptions := &types.ListOptions{}
	listCmd := &flaeg.Command{
		Name:                  "list",
		Description:           "Display all current PRs.",
		Config:                listOptions,
		DefaultPointersConfig: &types.ListOptions{},
		Run: safe(func() error {
			return cmd.List(listOptions)
		}),
	}
	return listCmd
}

func createVersion() *flaeg.Command {
	versionCmd := &flaeg.Command{
		Name:                  "version",
		Description:           "Display the version.",
		Config:                &types.NoOption{},
		DefaultPointersConfig: &types.NoOption{},
		Run: func() error {
			meta.DisplayVersion()
			return nil
		},
	}
	return versionCmd
}

func requirePRNumbers(numbers types.PRNumbers, action string) error {
	if len(numbers) == 0 {
		return fmt.Errorf("you must provide a PR number. ex: 'prm %s -n 1235'", action)
	}
	return nil
}

func safe(fn func() error) func() error {
	return func() error {
		_, err := config.Get()
		if err != nil {
			err = initProject()
			if err != nil {
				return err
			}
		}

		return fn()
	}
}

func initProject() error {
	// Get all remotes
	remotes, err := local.GetRemotes()
	if err != nil {
		return err
	}

	// get global configuration
	confs, err := config.ReadFile()
	if err != nil {
		return err
	}

	repoDir, err := local.GetGitRepoRoot()
	if err != nil {
		return err
	}

	var remoteName string
	if len(remotes) == 1 {
		remoteName = remotes[0].Name
	} else {
		remoteName, err = choose.GitRemote(remotes)
		if err != nil {
			return err
		}
		if len(remoteName) == 0 || remoteName == choose.ExitLabel {
			return errors.New("no remote chosen: exit")
		}
	}

	confs = append(confs, config.Configuration{
		Directory:  repoDir,
		BaseRemote: remoteName,
	})

	return config.Save(confs)
}
