package waldo

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type GitInfo struct {
	access GitAccess
	branch string
	commit string
}

//-----------------------------------------------------------------------------

type GitAccess int

const (
	Ok GitAccess = iota + 1 // MUST be first
	NoGitCommandFound
	NotGitRepository
)

func (ga GitAccess) String() string {
	return [...]string{
		"ok",
		"noGitCommandFound",
		"notGitRepository"}[ga-1]
}

//-----------------------------------------------------------------------------

func InferGitInfo(skipCount int) *GitInfo {
	access := Ok
	branch := ""
	commit := ""

	if !isGitInstalled() {
		access = NoGitCommandFound
	} else if !hasGitRepository() {
		access = NotGitRepository
	} else {
		commit = inferGitCommit(skipCount)
		branch = inferGitBranch(commit)
	}

	return &GitInfo{
		access: access,
		branch: branch,
		commit: commit}
}

//-----------------------------------------------------------------------------

func (gi *GitInfo) Access() GitAccess {
	return gi.access
}

func (gi *GitInfo) Branch() string {
	return gi.branch
}

func (gi *GitInfo) Commit() string {
	return gi.commit
}

//-----------------------------------------------------------------------------

func inferGitBranch(commit string) string {
	if len(commit) > 0 {
		name, _, err := run("git", "name-rev", "--always", "--name-only", commit)

		if err == nil {
			if strings.HasPrefix(name, "remotes/origin/") {
				name = name[len("remotes/origin/"):]
			}

			if name != "HEAD" {
				return name
			}
		}
	}

	name, _, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")

	if err == nil && name != "HEAD" {
		return name
	}

	return ""
}

func inferGitCommit(skipCount int) string {
	skip := fmt.Sprintf("--skip=%d", skipCount)

	hash, _, err := run("git", "log", "--format=%H", skip, "-1")

	if err != nil {
		return ""
	}

	return hash
}

func hasGitRepository() bool {
	_, _, err := run("git", "rev-parse")

	return err == nil
}

func isGitInstalled() bool {
	var name string

	if runtime.GOOS == "windows" {
		name = "git.exe"
	} else {
		name = "git"
	}

	_, err := exec.LookPath(name)

	return err == nil
}
