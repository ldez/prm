---
weight: 10
title: API Reference
---

# 🌱 Introduction

If you are a maintainer of Open Source Software, you need to review a lot of PR, this tool is made for you.

With the GitHub feature ["repository maintainer permissions on existing pull requests"](https://help.github.com/articles/allowing-changes-to-a-pull-request-branch-created-from-a-fork/), now we can edit real PR branch.
This tool allow to easily manage PR branches and remotes.

[![prm](https://asciinema.org/a/176222.png)](https://asciinema.org/a/176222)

💼 Features:

* Checkout a PR (interactively or by its number)
* Remove a PR (interactively or by its number)
* Remove all "checkouted" PRs for a project.
* Push on a PR.
* Display all "checkouted" PR (for a project or for all projects)
* Manage all your repositories.
* Save your configuration: `config/prm` (or `~/.prm` on Windows)
* Only works with GitHub.


# 💫 Checkout

## Interactive (Remote)

```bash
prm
# item "checkout"

# or #

prm checkout

# or #

prm c
```

* Display the last 50 PRs from GitHub.
* Add the user git remote named with the user login.
* Checkout the PR branch named like that: `<PR_NUMBER>--<BRANCH_NAME>`  
ex: `1234--myBranch`

## Interactive (Local)

```bash
prm
# item "list"
```

* Choose a PR between all "local" PRs.
* Checkout the PR branch named like that: `<PR_NUMBER>--<BRANCH_NAME>`  
ex: `1234--myBranch`

## By Number

```bash
prm checkout 1234

# or #

prm c 1234
```

* Add the user git remote named with the user login.
* Checkout the PR branch named like that: `<PR_NUMBER>--<BRANCH_NAME>`  
ex: `1234--myBranch`


# 💫 Remove

## Interactive

Only for the current project.

```bash
prm
# item "remove"

# or #

prm rm
```

* Display all "local" PRs.
* Remove by one or remove all.

## By Number

```bash
prm rm 1234

# or #

prm remove 1234
```

* Remove the local branch.
* Remove the user git remote if necessary.

## All

Only for the current project. (not all PR for all your projects)

```bash
prm rm --all
```

* Remove all PR related local branches.
* Remove all PR related git remote.

<aside class="notice">
It can be also done interactively with the item "remove".
</aside>

# 💫 Push

```bash
prm push
```

* Push to the PR related branch.
* Detect the PR number from the branch name.

## Push Force

```bash
prm pf
```

* Push force the PR related branch.
* Detect the PR number from the branch name.

# 💫 List

```bash
# display local branches related to PR. (current project only)
prm list

# display local branches related to PR. (all projects)
prm list --all
```

Display local branches related to PR for:

* current project
* all projects

# 💫 Clone

```bash
# clone and fork (if needed) a repository.
prm clone git@github.com:user/repo.git

# clone (don't create fork on GitHub and add a fork to remote.)
prm clone -n git@github.com:user/repo.git

# clone and fork (if needed) a repository, and use the username as root directory.
prm clone -r git@github.com:user/repo.git
```

Clone a repository:
* create a fork if needed
* add the remotes to the repository

# 💫 Help

```bash
prm -h

# or #

prm <command> -h
```

Display PRM help


# 🔒 Private Repositories

If you need to use `prm` for a private repository:

Create a [Github Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

## PRM_GITHUB_TOKEN

```bash
export PRM_GITHUB_TOKEN=xxxxxxx
```

Set the environment variable `PRM_GITHUB_TOKEN` with this token's value.

<aside class="notice">
You must replace `xxxxxxx` with your GitHub token.
</aside>

## PRM_GITHUB_TOKEN_FILE

Set the environment variable `PRM_GITHUB_TOKEN_FILE` with a path to file that contains this token's value.

```bash
export PRM_GITHUB_TOKEN_FILE=/path/to/my/token/secret.txt
```

<aside class="notice">
Important - `/path/to/my/token/secret.txt` contains only the value of the token.
</aside>


# 🔒 GitHub Enterprise

## PRM_GITHUB_API_BASE_URL

To use `prm` with GitHub Enterprise,

```bash
export PRM_GITHUB_API_BASE_URL=https://example.com/api/v3
```

Set the environment variable `PRM_GITHUB_API_BASE_URL` with the URL of the domain endpoint of your GitHub Enterprise instance.

<aside class="notice">
You must replace `https://example.com/api/v3` with the URL of the domain endpoint.
</aside>

## PRM_GITHUB_API_BASE_URL_FILE

Set the environment variable `PRM_GITHUB_API_BASE_URL_FILE` with a path to file that contains the URL of the domain endpoint of your GitHub Enterprise instance.

```bash
export PRM_GITHUB_API_BASE_URL_FILE=/path/to/my/token/secret.txt
```

<aside class="notice">
Important - `/path/to/my/token/secret.txt` contains only the URL of the domain endpoint.
</aside>


# 📦 How to Install

## Linux

### From Package Manager

> [ArchLinux (AUR)](https://aur.archlinux.org/packages/prm/)

```bash
yay -S prm
```

You can use a package manager:

* [ArchLinux (AUR)](https://aur.archlinux.org/packages/prm/)

### From Binaries

You can use pre-compiled binaries:

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/prm/releases)
* Unzip the archive.
* Add `prm` in your `PATH`.

## MacOS

### From Package Manager

> [Homebrew Taps](https://github.com/ldez/homebrew-tap)

```bash
brew tap ldez/tap
brew update
brew install prm
```

You can use a package manager:

* [Homebrew Taps](https://github.com/ldez/homebrew-tap)

### From Binaries

You can use pre-compiled binaries:

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/prm/releases)
* Unzip the archive.
* Add `prm` in your `PATH`.

## Windows

### From Package Manager

> [Scoop main bucket](https://github.com/lukesampson/scoop)

```bash
scoop install prm
```

> [Scoop Bucket](https://github.com/ldez/scoop-bucket)

```bash
scoop bucket add prm https://github.com/ldez/scoop-bucket.git
scoop install prm
```

You can use a package manager:

* [Scoop main bucket](https://github.com/lukesampson/scoop)
* [Scoop Bucket](https://github.com/ldez/scoop-bucket)

### From Binaries

You can use pre-compiled binaries:

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/prm/releases)
* Unzip the archive.
* Add `prm` in your `PATH`.

## From Sources

```bash
go get -u github.com/ldez/prm
```
