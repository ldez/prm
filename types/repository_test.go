package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	testCases := []struct {
		desc     string
		url      string
		expected *Repository
	}{
		{
			desc: "HTTPS",
			url:  "https://github.com/ldez/prm.git",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "SSH",
			url:  "git@github.com:containous/traefik.git",
			expected: &Repository{
				Owner: "containous",
				Name:  "traefik",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			repository, err := newRepository(test.url)

			require.NoError(t, err)
			assert.Equal(t, test.expected, repository)
		})
	}
}

func TestNewRepository_should_fail_when_invalid_URL(t *testing.T) {
	url := "https://github.com/ldez/prm"

	_, err := newRepository(url)

	require.Error(t, err)
}
