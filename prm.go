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
	"github.com/ldez/prm/meta"
	"github.com/ldez/prm/types"
	"github.com/ogier/pflag"
)

func main() {
	rootConfig := &types.RootOptions{}
	rootCmd := &flaeg.Command{
		Name:                  "prm",
		Description:           "PRM - The Pull Request Manager.",
		Config:                rootConfig,
		DefaultPointersConfig: &types.RootOptions{},
		Run: func() error {
			return cmd.Switch(&types.ListOptions{All: rootConfig.All})
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])

	flag.AddParser(reflect.TypeOf(types.PRNumbers{}), &types.PRNumbers{})

	// Checkout

	checkoutOptions := &types.CheckoutOptions{}

	checkoutCmd := &flaeg.Command{
		Name:                  "c",
		Description:           "Checkout a PR (create a local branch and add remote).",
		Config:                checkoutOptions,
		DefaultPointersConfig: &types.CheckoutOptions{},
	}
	checkoutCmd.Run = func() error {
		err := requirePRNumber(checkoutOptions.Number, checkoutCmd.Name)
		if err != nil {
			return err
		}
		return cmd.Checkout(checkoutOptions)
	}

	flag.AddCommand(checkoutCmd)

	// Get

	getOptions := &types.NoOption{}

	getCmd := &flaeg.Command{
		Name:                  "g",
		Description:           "Get remote PRs. (The last 25 PRs)",
		Config:                getOptions,
		DefaultPointersConfig: &types.NoOption{},
		Run: func() error {
			return cmd.Checkout(&types.CheckoutOptions{})
		},
	}

	flag.AddCommand(getCmd)

	// Remove

	removeOptions := &types.RemoveOptions{}

	removeCmd := &flaeg.Command{
		Name:                  "rm",
		Description:           "Remove one or more PRs from the current local repository.",
		Config:                removeOptions,
		DefaultPointersConfig: &types.RemoveOptions{},
	}
	removeCmd.Run = removeRun(removeCmd.Name, removeOptions)

	flag.AddCommand(removeCmd)

	// Push Force

	pushForceOptions := &types.PushOptions{Force: true}

	pushForceCmd := &flaeg.Command{
		Name:                  "pf",
		Description:           "Push force to the PR branch.",
		Config:                pushForceOptions,
		DefaultPointersConfig: &types.PushOptions{},
	}
	pushForceCmd.Run = func() error {
		return cmd.Push(pushForceOptions)
	}

	flag.AddCommand(pushForceCmd)

	// Push

	pushOptions := &types.PushOptions{}

	pushCmd := &flaeg.Command{
		Name:                  "push",
		Description:           "Push to the PR branch.",
		Config:                pushOptions,
		DefaultPointersConfig: &types.PushOptions{},
	}
	pushCmd.Run = func() error {
		return cmd.Push(pushOptions)
	}

	flag.AddCommand(pushCmd)

	// Pull

	pullOptions := &types.PullOptions{}

	pullCmd := &flaeg.Command{
		Name:                  "pull",
		Description:           "Pull to the PR branch.",
		Config:                pullOptions,
		DefaultPointersConfig: &types.PullOptions{},
	}
	pullCmd.Run = func() error {
		return cmd.Pull(pullOptions)
	}

	flag.AddCommand(pullCmd)

	// List

	listOptions := &types.ListOptions{}

	listCmd := &flaeg.Command{
		Name:                  "list",
		Description:           "Display all current PRs.",
		Config:                listOptions,
		DefaultPointersConfig: &types.ListOptions{},
		Run: func() error {
			return cmd.List(listOptions)
		},
	}

	flag.AddCommand(listCmd)

	// version

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

	flag.AddCommand(versionCmd)

	// Run command
	err := flag.Run()
	if err != nil && err != pflag.ErrHelp {
		log.Printf("Error: %v\n", err)
	}
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

			number, err := choose.PullRequest(conf.PullRequests)
			if err != nil || number <= choose.ExitValue {
				return err
			}

			if number == choose.AllValue {
				removeOptions.All = true
			} else {
				removeOptions.Numbers = append(removeOptions.Numbers, number)
			}

		} else {
			err := requirePRNumbers(removeOptions.Numbers, action)
			if err != nil {
				return err
			}
		}

		return cmd.Remove(removeOptions)
	}
}

func requirePRNumber(number int, action string) error {
	if number <= 0 {
		return fmt.Errorf("you must provide a PR number. ex: 'prm %s -n 1235'", action)
	}
	return nil
}

func requirePRNumbers(numbers types.PRNumbers, action string) error {
	if len(numbers) == 0 {
		return fmt.Errorf("you must provide a PR number. ex: 'prm %s -n 1235'", action)
	}
	return nil
}
