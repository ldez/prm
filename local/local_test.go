package local

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parsePRNumber_should_return_number_when_branch_name_respect_pattern(t *testing.T) {
	number, err := parsePRNumber("1234--branch")

	require.NoError(t, err)

	assert.Equal(t, 1234, number)
}

func Test_parsePRNumber(t *testing.T) {
	testCases := []struct {
		desc       string
		branchName string
	}{
		{
			desc:       "should fail when branch name don't respect pattern",
			branchName: "1234-branch",
		},
		{
			desc:       "should fail when branch name don't contain number",
			branchName: "xxxx-branch",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			_, err := parsePRNumber(test.branchName)
			require.Error(t, err)
		})
	}
}

func Test_parseRemotes(t *testing.T) {
	output := `
origin	git@github.com:ldez/traefik.git (fetch)
origin	git@github.com:ldez/traefik.git (push)
upstream	git@github.com:containous/traefik.git (fetch)
upstream	git@github.com:containous/traefik.git (push)
`

	remotes := parseRemotes(output)

	assert.Len(t, remotes, 2, "Wrong number of remotes: %v", remotes)
	assert.Equal(t, "origin", remotes[0].Name)
	assert.NotEmpty(t, remotes[0].URL)
	assert.Equal(t, "upstream", remotes[1].Name)
	assert.NotEmpty(t, remotes[1].URL)
}

func Test_parseRemotes_mix_spaces(t *testing.T) {
	output := `
origin  git@github.com:ldez/test.git (fetch)
origin  git@github.com:ldez/test.git (push)
upstream        https://github.com/ldez/test.git (fetch)
upstream	https://github.com/ldez/test.git (push)
`

	remotes := parseRemotes(output)

	assert.Len(t, remotes, 2, "Wrong number of remotes: %v", remotes)
	assert.Equal(t, "origin", remotes[0].Name)
	assert.NotEmpty(t, remotes[0].URL)
	assert.Equal(t, "upstream", remotes[1].Name)
	assert.NotEmpty(t, remotes[1].URL)
}

func Test_parseRemotes_empty_output(t *testing.T) {
	output := ``

	remotes := parseRemotes(output)

	assert.Len(t, remotes, 0, "Wrong number of remotes: %v", remotes)
}

func TestRemotes_Find_should_return_remote_when_exists(t *testing.T) {
	remoteName := "ccc"
	remotes := Remotes{
		{Name: "aaa", URL: "git@github.com:hubert/aaa.git"},
		{Name: "bbb", URL: "git@github.com:robert/aaa.git"},
		{Name: "ccc", URL: "git@github.com:norbert/aaa.git"},
		{Name: "ddd", URL: "git@github.com:gilbert/aaa.git"},
	}

	remote, err := remotes.Find(remoteName)

	require.NoError(t, err)
	assert.Equal(t, "ccc", remote.Name)
}

func TestRemotes_Find_should_thrown_error_when_not_exists(t *testing.T) {
	remoteName := "fff"
	remotes := Remotes{
		{Name: "aaa", URL: "git@github.com:hubert/aaa.git"},
		{Name: "bbb", URL: "git@github.com:robert/aaa.git"},
		{Name: "ccc", URL: "git@github.com:norbert/aaa.git"},
		{Name: "ddd", URL: "git@github.com:gilbert/aaa.git"},
	}

	_, err := remotes.Find(remoteName)

	require.Error(t, err)
}
