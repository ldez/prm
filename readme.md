# PRM - Pull Request Manager

[![Build Status](https://travis-ci.org/ldez/prm.svg?branch=master)](https://travis-ci.org/ldez/prm)
[![Go Report Card](https://goreportcard.com/badge/github.com/ldez/prm)](https://goreportcard.com/report/github.com/ldez/prm)


If you are a maintainer of Open Source Software, you need to review a lot of PR, this tool is made for you.

With the GitHub feature ["repository maintainer permissions on existing pull requests"](https://help.github.com/articles/allowing-changes-to-a-pull-request-branch-created-from-a-fork/), now we can edit real PR branch.
This tool allow to easily manage PR branches and remotes.

[![prm](https://asciinema.org/a/176222.png)](https://asciinema.org/a/176222)

:briefcase: Features:

* Checkout a PR (interactively or by its number)
* Remove a PR (interactively or by its number)
* Remove all "checkouted" PRs for a project.
* Push on a PR.
* Display all "checkouted" PR (for a project or for all projects)
* Manage all your repositories.
* Save your configuration: `config/prm` (or `~/.prm` on Windows)
* Only works with GitHub.

---

:package: **How to install:**
**[From binaries](#from-binaries)** -
**[From a package manager](#from-a-package-manager)** -
**[From sources](#from-sources)** 

:bulb: **How to use:**
**[Checkout](#checkout)** -
**[Remove](#remove)** -
**[Push](#push)** -
**[Push Force](#push-force)** -
**[List](#list)** -
**[Help](#help)**

---

## How to Install

### From Binaries

* To get the binary just download the latest release for your OS/Arch from [the releases page](https://github.com/ldez/prm/releases)
* Unzip the archive.
* Add `prm` in your `PATH`.

Available for: Linux, MacOS, Windows, FreeBSD, OpenBSD.

### From a Package Manager

- [ArchLinux (AUR)](https://aur.archlinux.org/packages/prm/)
```bash
yay -S prm
```

- [Homebrew Taps](https://github.com/ldez/homebrew-tap)
```bash
brew tap ldez/tap
brew update
brew install prm
```

- [Scoop main bucket](https://github.com/lukesampson/scoop)
```bash
scoop install prm
```

- [Scoop Bucket](https://github.com/ldez/scoop-bucket)
```bash
scoop bucket add prm https://github.com/ldez/scoop-bucket.git
scoop install prm
```

### From Sources

```bash
go get -u github.com/ldez/prm
```

### Configuration for Private Repositories

If you need to use `prm` for a private repository:

* Create a [Github Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)
* Set the environment variable `PRM_GITHUB_TOKEN` with this token's value:
```bash
export PRM_GITHUB_TOKEN=xxxxxxx
```
* Or set the environment variable `PRM_GITHUB_TOKEN_FILE` with a path to file that contains this token's value.
```bash
# /path/to/my/token/secret.txt contains only the value of the token.
export PRM_GITHUB_TOKEN_FILE=/path/to/my/token/secret.txt
```


## How to Use

### Checkout

#### Interactive (Remote)

```bash
prm
# item "checkout"
```

* Display the last 50 PRs from GitHub.
* Add the user git remote named with the user login.
* Checkout the PR branch named like that: `<PR_NUMBER>--<BRANCH_NAME>`
ex: `1234\--myBranch`

#### Interactive (Local)

```bash
prm
# item "list"
```

* Choose a PR between all "local" PRs.
* Checkout the PR branch named like that: `<PR_NUMBER>--<BRANCH_NAME>`
ex: `1234\--myBranch`

#### By Number

```bash
prm c -n 1234
# or
prm c --number=1234
```

* Add the user git remote named with the user login.
* Checkout the PR branch named like that: `<PR_NUMBER>--<BRANCH_NAME>`
ex: `1234\--myBranch`

### Remove

#### Interactive

Only for the current project.

```bash
prm rm

# or
prm
# item "remove"
```

* Display all "local" PRs.
* Remove by one or remove all.

#### By Number

```bash
prm rm -n 1234
# or
prm rm --number=1234
```

* Remove the local branch.
* Remove the user git remote if necessary.

#### All

Only for the current project. (not all PR for all your projects)

```bash
prm rm --all
```

* Remove all PR related local branches.
* Remove all PR related git remote.

:star: It can be also done interactively with the item "remove".

### Push

```bash
prm push
```

* Push to the PR related branch.
* Detect the PR number from the branch name.

### Push Force

```bash
prm pf
```

* Push force the PR related branch.
* Detect the PR number from the branch name.

### List

```bash
# display local branches related to PR. (current project only)
prm list

# display local branches related to PR. (all projects)
prm list --all
```

* Display local branches related to PR for:
  * current project
  * all projects

### Help

```bash
prm -h
```
