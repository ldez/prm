package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_splitUserRepo(t *testing.T) {
	testCases := []string{
		"git@github.com:ldez/prm.git",
		"git@github.com:ldez/prm",
		"https://github.com/ldez/prm.git",
		"https://github.com/ldez/prm",
	}

	for _, test := range testCases {
		user, repo, err := splitUserRepo(test)
		require.NoError(t, err)

		assert.Equal(t, "ldez", user)
		assert.Equal(t, "prm", repo)
	}
}

func Test_searchFork(t *testing.T) {
	t.Skip("e2e")

	ctx := t.Context()

	cl := newCloner(ctx)

	// Special cases when the repository name of the fork is not the same as the parent repository name.

	// vdemeester/docker-cli forked from docker/cli
	// vdemeester/openshift-release forked from openshift/release

	repo, err := cl.searchFork(ctx, "vdemeester", "docker", "cli")
	require.NoError(t, err)

	require.NotNil(t, repo.Parent)
	assert.Equal(t, "docker/cli", repo.GetParent().GetFullName())

	repo, err = cl.searchFork(ctx, "vdemeester", "openshift", "release")
	require.NoError(t, err)

	require.NotNil(t, repo.Parent)
	assert.Equal(t, "openshift/release", repo.GetParent().GetFullName())
}
