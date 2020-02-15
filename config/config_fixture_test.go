package config

import "github.com/ldez/prm/v3/types"

const (
	defaultTestBaseRemote = "remoteA"
	defaultTestDirectory  = "/my/git/directory/aaa"
)

func aConfiguration(builders ...func(*Configuration)) Configuration {
	conf := &Configuration{
		BaseRemote:   defaultTestBaseRemote,
		Directory:    defaultTestDirectory,
		PullRequests: make(map[string][]types.PullRequest),
	}

	for _, builder := range builders {
		builder(conf)
	}

	return *conf
}

func directoryB(conf *Configuration) {
	conf.Directory = "/my/git/directory/bbb"
}

func directoryC(conf *Configuration) {
	conf.Directory = "/my/git/directory/ccc"
}

func withPullRequest(owner string, builders ...func(*types.PullRequest)) func(*Configuration) {
	return func(conf *Configuration) {
		pr := &types.PullRequest{
			Owner: owner,
		}
		for _, builder := range builders {
			builder(pr)
		}
		conf.PullRequests[owner] = append(conf.PullRequests[owner], *pr)
	}
}

func aPullRequest(owner string, builders ...func(pr *types.PullRequest)) types.PullRequest {
	pr := &types.PullRequest{
		Owner: owner,
	}
	for _, builder := range builders {
		builder(pr)
	}

	return *pr
}

func branchA(pr *types.PullRequest) {
	pr.BranchName = "branchA"
}

func branchB(pr *types.PullRequest) {
	pr.BranchName = "branchB"
}

func withNumber(number int) func(*types.PullRequest) {
	return func(pr *types.PullRequest) {
		pr.Number = number
	}
}
