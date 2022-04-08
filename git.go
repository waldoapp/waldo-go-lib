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

func removeDuplicate(strings []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strings {
		if _, ok := seen[s]; !ok {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func refNameToBranchName(refName string) string {
	branchName := strings.TrimSpace(refName)

	if strings.HasPrefix(branchName, "refs/heads/") {
		branchName = strings.TrimPrefix(branchName, "refs/heads/")
	} else if strings.HasPrefix(branchName, "refs/remotes/") {
		branchName = strings.TrimPrefix(branchName, "refs/remotes/")

		// Remove the remote name
		slash := strings.Index(branchName, "/")
		if slash == -1 {
			return ""
		}
		branchName = branchName[slash+1:]
	} else {
		return ""
	}

	if branchName == "HEAD" {
		return ""
	}

	return branchName
}

func getBranchNamesFromGitForeachRefResults(results string) []string {
	lines := strings.Split(results, "\n")

	var branchNames []string

	for _, line := range lines {
		branchName := refNameToBranchName(line)

		if len(branchName) > 0 {
			branchNames = append(branchNames, branchName)
		}
	}

	return removeDuplicate(branchNames)
}

func inferGitBranchFromForEachRef(commit string) string {
	stdout, _, err := run("git", "for-each-ref", fmt.Sprintf("--points-at=%s", commit), "--format=%(refname)")

	if err == nil {
		branchNames := getBranchNamesFromGitForeachRefResults(stdout)
		if len(branchNames) > 0 {
			// We don't know which branch is the right one, so just return the first one
			return branchNames[0]
		}
	}

	return ""
}

func inferGitBranchFromNameRev(commit string) string {
	name, _, err := run("git", "name-rev", "--always", "--name-only", commit)

	if err == nil {
		return refNameToBranchName(name)
	}

	return ""
}

func inferGitBranchFromRevParse() string {
	name, _, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")

	if err == nil && name != "HEAD" {
		return name
	}

	return ""
}

func inferGitBranch(commit string) string {
	if len(commit) > 0 {
		fromForeachRev := inferGitBranchFromForEachRef(commit)
		if fromForeachRev != "" {
			return fromForeachRev
		}

		fromNameRev := inferGitBranchFromNameRev(commit)
		if fromNameRev != "" {
			return fromNameRev
		}
	}

	return inferGitBranchFromRevParse()
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
