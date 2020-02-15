package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/ldez/prm/v3/choose"
	"github.com/ldez/prm/v3/cmd"
	"github.com/ldez/prm/v3/config"
	"github.com/ldez/prm/v3/local"
	"github.com/ldez/prm/v3/meta"
	"github.com/ldez/prm/v3/types"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := createRootCmd()
	rootCmd.AddCommand(createCheckoutCmd())
	rootCmd.AddCommand(createRemoveCmd())
	rootCmd.AddCommand(createPullCmd())
	rootCmd.AddCommand(createPushCmd())
	rootCmd.AddCommand(createPushForceCmd())
	rootCmd.AddCommand(createListCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "prm",
		Short:   "PRM - The Pull Request Manager.",
		Long:    `PRM - The Pull Request Manager.`,
		Version: meta.GetVersion(),
		PreRunE: safe,
		RunE: func(_ *cobra.Command, _ []string) error {
			return rootRun()
		},
	}
}

func createCheckoutCmd() *cobra.Command {
	checkoutCfg := types.CheckoutOptions{}

	return &cobra.Command{
		Use:     "checkout [PR number]",
		Aliases: []string{"c"},
		Short:   "Checkout a PR (create a local branch and add remote).",
		Long:    "Checkout a PR (create a local branch and add remote).",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
				return err
			}

			if len(args) == 1 {
				val, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("argument must be a valid number: %w", err)
				}
				checkoutCfg.Number = val
			}

			return nil
		},
		PreRunE: safe,
		RunE: func(_ *cobra.Command, args []string) error {
			if checkoutCfg.Number != 0 {
				return cmd.Checkout(&checkoutCfg)
			}

			conf, err := config.Get()
			if err != nil {
				return err
			}
			return cmd.InteractiveCheckout(conf)
		},
		Example: `  $ prm checkout
  $ prm checkout 1234
  $ prm c
  $ prm c 1234`,
	}
}

func createRemoveCmd() *cobra.Command {
	removeCfg := types.RemoveOptions{}

	removeCmd := &cobra.Command{
		Use:     "rm [PR numbers]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more PRs from the current local repository.",
		Long:    "Remove one or more PRs from the current local repository.",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			var values []int
			for i, arg := range args {
				val, err := strconv.Atoi(arg)
				if err != nil {
					return fmt.Errorf("argument %d must be a valid number: %w", i, err)
				}
				values = append(values, val)
			}

			removeCfg.Numbers = values

			return nil
		},
		PreRunE: safe,
		RunE: func(_ *cobra.Command, args []string) error {
			return removeRun(&removeCfg)
		},
		Example: `  $ prm rm
  $ prm rm 1234
  $ prm rm 1234 4567
  $ prm remove
  $ prm remove 1234
  $ prm remove 1234 4567`,
	}

	removeCmd.Flags().BoolVarP(&removeCfg.All, "all", "a", false, "All pull requests.")

	return removeCmd
}

func createPullCmd() *cobra.Command {
	pullCfg := types.PullOptions{}

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull to the PR branch.",
		Long:  "Pull to the PR branch.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.Pull(&pullCfg)
		},
	}

	pullCmd.Flags().BoolVarP(&pullCfg.Force, "force", "f", false, "Force the pull.")

	return pullCmd
}

func createPushCmd() *cobra.Command {
	pushCfg := types.PushOptions{}

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push to the PR branch.",
		Long:  "Push to the PR branch.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.Push(&pushCfg)
		},
	}

	pushCmd.Flags().BoolVarP(&pushCfg.Force, "force", "f", false, "Force the push.")

	return pushCmd
}

func createPushForceCmd() *cobra.Command {
	pushForceCfg := types.PushOptions{Force: true}

	return &cobra.Command{
		Use:   "pf",
		Short: "Push force to the PR branch.",
		Long:  "Push force to the PR branch.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.Push(&pushForceCfg)
		},
	}
}

func createListCmd() *cobra.Command {
	listCfg := types.ListOptions{}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Display all current PRs.",
		Long:  "Display all current PRs.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.List(&listCfg)
		},
	}

	listCmd.Flags().BoolVarP(&listCfg.All, "all", "a", false, "All pull requests.")

	return listCmd
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

func removeRun(removeOptions *types.RemoveOptions) error {
	if removeOptions.All {
		return cmd.Remove(removeOptions)
	}

	if len(removeOptions.Numbers) == 0 {
		conf, err := config.Get()
		if err != nil {
			return err
		}

		return cmd.InteractiveRemove(conf)
	}

	return cmd.Remove(removeOptions)
}

func safe(_ *cobra.Command, _ []string) error {
	_, err := config.Get()
	if err == nil {
		return nil
	}

	return initProject()
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
