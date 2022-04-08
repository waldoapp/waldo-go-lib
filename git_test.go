package waldo

import (
	"testing"
)

func TestGetBranchesEmptyString(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("")

	if names != nil {
		t.Errorf("Expected nil, got %v", names)
	}
}

func TestGetBranchesIgnoreTags(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/tags/foo")

	if names != nil {
		t.Errorf("Expected nil, got %v", names)
	}
}

func TestGetBranchesIgnoreHeads(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/remotes/origin/HEAD")

	if names != nil {
		t.Errorf("Expected nil, got %v", names)
	}
}

func TestGetBranchesLocalMaster(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/heads/master")

	if len(names) != 1 {
		t.Errorf("Expected 1 element, got %v", names)
	}

	if names[0] != "master" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}
}

func TestGetBranchesOriginMaster(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/remotes/origin/master")

	if len(names) != 1 {
		t.Errorf("Expected 1 element, got %v", names)
	}

	if names[0] != "master" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}
}

func TestGetBranchesCanHaveMultiple(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/remotes/origin/master\nrefs/heads/foo\nrefs/remotes/origin/bar")

	if len(names) != 3 {
		t.Errorf("Expected 1 element, got %v", names)
	}

	if names[0] != "master" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}

	if names[1] != "foo" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}

	if names[2] != "bar" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}
}

func TestGetBranchesDeduplicates(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/remotes/origin/master\nrefs/heads/master")

	if len(names) != 1 {
		t.Errorf("Expected 1 element, got %v", names)
	}

	if names[0] != "master" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}
}

func TestGetBranchesFull(t *testing.T) {
	names := getBranchNamesFromGitForeachRefResults("refs/heads/master\nrefs/remotes/origin/HEAD\nrefs/remotes/origin/master\nrefs/tags/foo")

	if len(names) != 1 {
		t.Errorf("Expected 1 element, got %v", names)
	}

	if names[0] != "master" {
		t.Errorf("Expected [\"master\"], got %v", names)
	}
}
