package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	giturls "github.com/whilp/git-urls"
)

func Test_splitUserRepo(t *testing.T) {
	testCases := []string{
		"git@github.com:ldez/prm.git",
		"git@github.com:ldez/prm",
		"https://github.com/ldez/prm.git",
		"https://github.com/ldez/prm",
	}

	for _, test := range testCases {
		u, err := giturls.Parse(test)
		require.NoError(t, err)

		user, repo, err := splitUserRepo(u)
		require.NoError(t, err)

		assert.Equal(t, "ldez", user)
		assert.Equal(t, "prm", repo)
	}
}
