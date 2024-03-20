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
			desc: "HTTPS without suffix",
			url:  "https://github.com/ldez/prm",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "HTTPS without suffix ending with /",
			url:  "https://github.com/ldez/prm/",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "SSH",
			url:  "git@github.com:ldez/prm.git",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "SSH without suffix",
			url:  "git@github.com:ldez/prm",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "GitHub Enterprise: HTTPS",
			url:  "https://github.mycompany.com/ldez/prm.git",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "GitHub Enterprise: HTTPS without suffix",
			url:  "https://github.mycompany.com/ldez/prm",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "GitHub Enterprise: HTTPS without suffix ending with /",
			url:  "https://github.mycompany.com/ldez/prm/",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "GitHub Enterprise: SSH",
			url:  "git@github.mycompany.com:ldez/prm.git",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
		{
			desc: "GitHub Enterprise: SSH without suffix",
			url:  "git@github.mycompany.com:ldez/prm",
			expected: &Repository{
				Owner: "ldez",
				Name:  "prm",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			repository, err := newRepository(test.url)
			require.NoError(t, err)

			assert.Equal(t, test.expected, repository)
		})
	}
}
