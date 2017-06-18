package cmd

import (
	"testing"

	"github.com/ldez/prm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newRepository_https(t *testing.T) {
	url := "https://github.com/ldez/prm.git"

	repository, err := newRepository(url)

	require.NoError(t, err)
	assert.Equal(t, "prm", repository.Name)
	assert.Equal(t, "ldez", repository.Owner)
}

func Test_newRepository_ssh(t *testing.T) {
	url := "git@github.com:containous/traefik.git"

	repository, err := newRepository(url)

	require.NoError(t, err)
	assert.Equal(t, "traefik", repository.Name)
	assert.Equal(t, "containous", repository.Owner)
}

func Test_getRemotes(t *testing.T) {
	output := `
origin	git@github.com:ldez/traefik.git (fetch)
origin	git@github.com:ldez/traefik.git (push)
upstream	git@github.com:containous/traefik.git (fetch)
upstream	git@github.com:containous/traefik.git (push)
`

	remotes := getRemotes(output)

	assert.Len(t, remotes, 2, "Wrong number of remotes: %v", remotes)
	origin := remotes[0]
	assert.Equal(t, "origin", origin.Name)
	upstream := remotes[1]
	assert.Equal(t, "upstream", upstream.Name)
}

func Test_getRemotes_empty_output(t *testing.T) {
	output := ``

	remotes := getRemotes(output)

	assert.Len(t, remotes, 0, "Wrong number of remotes: %v", remotes)
}

func Test_findRemote_should_return_remote_when_exists(t *testing.T) {
	remoteName := "ccc"
	remotes := []types.Remote{
		{Name: "aaa", URL: "git@github.com:hubert/aaa.git"},
		{Name: "bbb", URL: "git@github.com:robert/aaa.git"},
		{Name: "ccc", URL: "git@github.com:norbert/aaa.git"},
		{Name: "ddd", URL: "git@github.com:gilbert/aaa.git"},
	}

	remote, err := findRemote(remotes, remoteName)

	require.NoError(t, err)
	assert.Equal(t, "ccc", remote.Name)
}

func Test_findRemote_should_thrown_error_when_not_exists(t *testing.T) {
	remoteName := "fff"
	remotes := []types.Remote{
		{Name: "aaa", URL: "git@github.com:hubert/aaa.git"},
		{Name: "bbb", URL: "git@github.com:robert/aaa.git"},
		{Name: "ccc", URL: "git@github.com:norbert/aaa.git"},
		{Name: "ddd", URL: "git@github.com:gilbert/aaa.git"},
	}

	_, err := findRemote(remotes, remoteName)

	require.Error(t, err)
}
