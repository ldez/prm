package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/containous/flaeg"
	"github.com/ldez/prm/cmd"
	"github.com/ldez/prm/meta"
	"github.com/ldez/prm/types"
)

func main() {
	emptyConfig := &types.NoOption{}
	rootCmd := &flaeg.Command{
		Name:                  "prm",
		Description:           "PRM - The Pull Request Manager.",
		Config:                emptyConfig,
		DefaultPointersConfig: &types.NoOption{},
		Run: func() error {
			return cmd.List(&types.ListOptions{})
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

	// Remove

	removeOptions := &types.RemoveOptions{}

	removeCmd := &flaeg.Command{
		Name:                  "rm",
		Description:           "Remove one or more PRs from the current local repository.",
		Config:                removeOptions,
		DefaultPointersConfig: &types.RemoveOptions{},
	}
	removeCmd.Run = func() error {
		err := requirePRNumbers(removeOptions.Numbers, removeCmd.Name)
		if !removeOptions.All && err != nil {
			return err
		}
		return cmd.Remove(removeOptions)
	}

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
	if err != nil && !strings.HasSuffix(err.Error(), "pflag: help requested") {
		log.Printf("Error: %v\n", err)
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
