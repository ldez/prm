package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfiguration_RemovePullRequest_should_remove_all_when_only_one_PR_for_a_owner(t *testing.T) {
	conf := aConfiguration(
		withPullRequest("hubert", withNumber(1), branchA),
		withPullRequest("hubert", withNumber(2), branchB),
		withPullRequest("robert", withNumber(3), branchA),
		withPullRequest("norbert", withNumber(4), branchA))

	pr := aPullRequest("norbert", withNumber(4), branchA)

	count := conf.RemovePullRequest(&pr)

	assert.Equalf(t, 0, count, "in configuration: %v", conf)
	assert.Emptyf(t, conf.PullRequests["norbert"], "in configuration: %v", conf)
}

func TestConfiguration_RemovePullRequest_should_remove_only_PR_when_multiple_PR_for_a_owner(t *testing.T) {
	conf := aConfiguration(
		withPullRequest("hubert", withNumber(1), branchA),
		withPullRequest("hubert", withNumber(2), branchB),
		withPullRequest("robert", withNumber(3), branchA),
		withPullRequest("norbert", withNumber(4), branchA))

	pr := aPullRequest("hubert", withNumber(1), branchA)

	count := conf.RemovePullRequest(&pr)

	assert.Equalf(t, 1, count, "in configuration: %v", conf)
	assert.Lenf(t, conf.PullRequests["hubert"], 1, "in configuration: %v", conf)
}

func TestConfiguration_FindPullRequests_should_return_pr_when_pr_exist(t *testing.T) {
	conf := aConfiguration(
		withPullRequest("hubert", withNumber(1), branchA),
		withPullRequest("hubert", withNumber(2), branchB),
		withPullRequest("robert", withNumber(3), branchA),
		withPullRequest("norbert", withNumber(4), branchA))

	number := 4

	pr, err := conf.FindPullRequests(number)

	require.NoErrorf(t, err, "in configuration: %v", conf)
	assert.Equalf(t, number, pr.Number, "in configuration: %v", conf)
	assert.Equalf(t, "norbert", pr.Owner, "in configuration: %v", conf)
}

func TestConfiguration_FindPullRequests_should_fail_when_pr_not_exist(t *testing.T) {
	conf := aConfiguration(
		withPullRequest("hubert", withNumber(1), branchA),
		withPullRequest("hubert", withNumber(2), branchB),
		withPullRequest("robert", withNumber(3), branchA),
		withPullRequest("norbert", withNumber(4), branchA))

	number := 5

	_, err := conf.FindPullRequests(number)

	require.Errorf(t, err, "in configuration: %v", conf)
}

func TestReadFile_should_return_empty_configuration_list_when_file_not_exist(t *testing.T) {
	dir := t.TempDir()

	// Mock GetPath function
	getPathFunc = func() (string, error) {
		return path.Join(dir, defaultFileName), nil
	}

	confs, err := ReadFile()

	require.NoError(t, err)
	assert.Len(t, confs, 0)
}

func TestReadFile_should_return_configuration_list_when_file_exist(t *testing.T) {
	// Mock GetPath function
	getPathFunc = func() (string, error) {
		return path.Join("fixture", "case01.json"), nil
	}

	confs, err := ReadFile()

	require.NoError(t, err)
	assert.Len(t, confs, 2)
}

func TestSave_should_save_configuration(t *testing.T) {
	dir := t.TempDir()

	// Mock GetPath function
	getPathFunc = func() (string, error) {
		return path.Join(dir, defaultFileName), nil
	}

	cfs := []Configuration{
		aConfiguration(directoryB, withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(directoryC, withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(withPullRequest("hubert", withNumber(2), branchB)),
	}

	err := Save(cfs)

	require.NoError(t, err)
}

func TestFind_should_find_configuration_when_configuration_exists(t *testing.T) {
	directory := defaultTestDirectory
	cfs := []Configuration{
		aConfiguration(directoryB, withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(directoryC, withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(withPullRequest("hubert", withNumber(2), branchB)),
	}

	conf, err := Find(cfs, directory)

	require.NoErrorf(t, err, "Error during find.")
	assert.Equalf(t, defaultTestDirectory, conf.Directory, "It's not the right configuration: %s", conf)
}

func TestFind_should_fail_when_configuration_not_exists(t *testing.T) {
	directory := defaultTestDirectory
	cfs := []Configuration{
		aConfiguration(directoryB, withPullRequest("hubert", withNumber(1), branchA)),
		aConfiguration(directoryC, withPullRequest("hubert", withNumber(1), branchA)),
	}

	_, err := Find(cfs, directory)

	require.Errorf(t, err, "Must thrown an error. Configuration: %v, Directory: %s", cfs, directory)
}
