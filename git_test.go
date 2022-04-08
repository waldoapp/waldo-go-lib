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

func TestNameRevToBranchNameEmpty(t *testing.T) {
	name := nameRevToBranchName("")
	expected := ""
	if name != expected {
		t.Errorf("Expected %s string, got %v", expected, name)
	}
}

func TestNameRevToBranchNameSimple(t *testing.T) {
	name := nameRevToBranchName("foo")
	expected := "foo"
	if name != expected {
		t.Errorf("Expected %s string, got %v", expected, name)
	}
}

func TestNameRevToBranchNameComplex(t *testing.T) {
	name := nameRevToBranchName("features/waldo/git-handling")
	expected := "features/waldo/git-handling"
	if name != expected {
		t.Errorf("Expected %s string, got %v", expected, name)
	}
}

func TestNameRevToBranchNameRemote(t *testing.T) {
	name := nameRevToBranchName("remotes/origin/foo")
	expected := "foo"
	if name != expected {
		t.Errorf("Expected %s string, got %v", expected, name)
	}
}

func TestNameRevToBranchNameTags(t *testing.T) {
	name := nameRevToBranchName("tags/bar")
	expected := ""
	if name != expected {
		t.Errorf("Expected %s string, got %v", expected, name)
	}
}

func TestNameRevToBranchNameHEAD(t *testing.T) {
	name := nameRevToBranchName("HEAD")
	expected := ""
	if name != expected {
		t.Errorf("Expected %s string, got %v", expected, name)
	}
}
