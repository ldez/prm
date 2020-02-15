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
