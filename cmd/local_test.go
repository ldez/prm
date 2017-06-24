package cmd

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

func Test_parsePRNumber_should_fail_when_branch_name_dont_respect_pattern(t *testing.T) {
	_, err := parsePRNumber("1234-branch")

	require.Error(t, err)
}

func Test_parsePRNumber_should_fail_when_branch_name_dont_contain_number(t *testing.T) {
	_, err := parsePRNumber("xxxx-branch")

	require.Error(t, err)
}
