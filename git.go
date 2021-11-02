package waldo

import (
	"fmt"
	"os/exec"
	"runtime"
)

func inferGit(skipCount int) (string, string, string) {
	if !isGitInstalled() {
		return "noGitCommandFound", "", ""
	}

	if !hasGitRepository() {
		return "notGitRepository", "", ""
	}

	gitCommit := inferGitCommit(skipCount)
	gitBranch := inferGitBranch(gitCommit)

	return "ok", gitBranch, gitCommit
}

func inferGitBranch(commit string) string {
	if len(commit) > 0 {
		name, _, err := run("git", "name-rev", "--refs=heads/*", "--name-only", commit)

		if err == nil && name != "HEAD" {
			return name
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
