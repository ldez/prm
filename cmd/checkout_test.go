package cmd

import (
	"testing"

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
