package waldo

import (
	"os"
	"strings"
)

type CIInfo struct {
	gitBranch string
	gitCommit string
	provider  CIProvider
	skipCount int
}

//-----------------------------------------------------------------------------

type CIProvider int

const (
	Unknown CIProvider = iota // MUST be first
	AppCenter
	AzureDevOps
	Bitrise
	CircleCI
	CodeBuild
	GitHubActions
	Jenkins
	TeamCity
	TravisCI
	XcodeCloud
)

func (cp CIProvider) String() string {
	return [...]string{
		"Unknown",
		"App Center",
		"Azure DevOps",
		"Bitrise",
		"CircleCI",
		"CodeBuild",
		"GitHub Actions",
		"Jenkins",
		"TeamCity",
		"Travis CI",
		"Xcode Cloud"}[cp]
}

//-----------------------------------------------------------------------------

func DetectCIInfo(fullInfo bool) *CIInfo {
	info := &CIInfo{provider: detectCIProvider()}

	if fullInfo {
		info.extractFullInfo()
	}

	return info
}

//-----------------------------------------------------------------------------

func (ci *CIInfo) GitBranch() string {
	return ci.gitBranch
}

func (ci *CIInfo) GitCommit() string {
	return ci.gitCommit
}

func (ci *CIInfo) Provider() CIProvider {
	return ci.provider
}

func (ci *CIInfo) SkipCount() int {
	return ci.skipCount
}

//-----------------------------------------------------------------------------

func (ci *CIInfo) extractFullInfo() {
	switch ci.provider {
	case AppCenter:
		ci.extractFullInfoFromAppCenter()

	case AzureDevOps:
		ci.extractFullInfoFromAzureDevOps()

	case Bitrise:
		ci.extractFullInfoFromBitrise()

	case CircleCI:
		ci.extractFullInfoFromCircleCI()

	case CodeBuild:
		ci.extractFullInfoFromCodeBuild()

	case GitHubActions:
		ci.extractFullInfoFromGitHubActions()

	case Jenkins:
		ci.extractFullInfoFromJenkins()

	case TeamCity:
		ci.extractFullInfoFromTeamCity()

	case TravisCI:
		ci.extractFullInfoFromTravisCI()

	case XcodeCloud:
		ci.extractFullInfoFromXcodeCloud()

	default:
		break
	}
}

func (ci *CIInfo) extractFullInfoFromAppCenter() {
	ci.gitBranch = os.Getenv("APPCENTER_BRANCH")
	ci.gitCommit = "" //os.Getenv("???") -- not currently supported?
}

func (ci *CIInfo) extractFullInfoFromAzureDevOps() {
	ci.gitBranch = os.Getenv("BUILD_SOURCEBRANCHNAME")
	ci.gitCommit = os.Getenv("BUILD_SOURCEVERSION")
}

func (ci *CIInfo) extractFullInfoFromBitrise() {
	ci.gitBranch = os.Getenv("BITRISE_GIT_BRANCH")
	ci.gitCommit = os.Getenv("BITRISE_GIT_COMMIT")
}

func (ci *CIInfo) extractFullInfoFromCircleCI() {
	ci.gitBranch = os.Getenv("CIRCLE_BRANCH")
	ci.gitCommit = os.Getenv("CIRCLE_SHA1")
}

func (ci *CIInfo) extractFullInfoFromCodeBuild() {
	trigger := os.Getenv("CODEBUILD_WEBHOOK_TRIGGER")

	if strings.HasPrefix(trigger, "branch/") {
		ci.gitBranch = strings.TrimPrefix(trigger, "branch/")
	} else {
		ci.gitBranch = ""
	}

	ci.gitCommit = os.Getenv("CODEBUILD_WEBHOOK_PREV_COMMIT")
}

func (ci *CIInfo) extractFullInfoFromGitHubActions() {
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	refType := os.Getenv("GITHUB_REF_TYPE")

	switch eventName {
	case "pull_request", "pull_request_target":
		if refType == "branch" {
			ci.gitBranch = os.Getenv("GITHUB_HEAD_REF")
		} else {
			ci.gitBranch = ""
		}

		//
		// The following environment variable must be set by us (most likely in
		// a custom action) to match the current value of
		// `github.event.pull_request.head.sha`:
		//
		ci.gitCommit = os.Getenv("GITHUB_EVENT_PULL_REQUEST_HEAD_SHA")

		ci.skipCount = 1

	case "push":
		if refType == "branch" {
			ci.gitBranch = os.Getenv("GITHUB_REF_NAME")
		} else {
			ci.gitBranch = ""
		}

		ci.gitCommit = os.Getenv("GITHUB_SHA")

	default:
		ci.gitBranch = ""
		ci.gitCommit = ""
	}
}

func (ci *CIInfo) extractFullInfoFromJenkins() {
	ci.gitBranch = "" //os.Getenv("???") -- not currently supported?
	ci.gitCommit = "" //os.Getenv("???") -- not currently supported?
}

func (ci *CIInfo) extractFullInfoFromTeamCity() {
	ci.gitBranch = "" //os.Getenv("???") -- not currently supported?
	ci.gitCommit = "" //os.Getenv("???") -- not currently supported?
}

func (ci *CIInfo) extractFullInfoFromTravisCI() {
	ci.gitBranch = os.Getenv("TRAVIS_BRANCH")
	ci.gitCommit = os.Getenv("TRAVIS_COMMIT")
}

func (ci *CIInfo) extractFullInfoFromXcodeCloud() {
	ci.gitBranch = os.Getenv("CI_BRANCH")
	ci.gitCommit = os.Getenv("CI_COMMIT")
}

//-----------------------------------------------------------------------------

func detectCIProvider() CIProvider {
	switch {
	case onAppCenter():
		return AppCenter

	case onAzureDevOps():
		return AzureDevOps

	case onBitrise():
		return Bitrise

	case onCircleCI():
		return CircleCI

	case onCodeBuild():
		return CodeBuild

	case onGitHubActions():
		return GitHubActions

	case onJenkins():
		return Jenkins

	case onTeamCity():
		return TeamCity

	case onTravisCI():
		return TravisCI

	case onXcodeCloud():
		return XcodeCloud

	default:
		return Unknown
	}
}

func onAppCenter() bool {
	return len(os.Getenv("APPCENTER_BUILD_ID")) > 0
}

func onAzureDevOps() bool {
	return len(os.Getenv("AGENT_ID")) > 0
}

func onBitrise() bool {
	return os.Getenv("BITRISE_IO") == "true"
}

func onCircleCI() bool {
	return os.Getenv("CIRCLECI") == "true"
}

func onCodeBuild() bool {
	return len(os.Getenv("CODEBUILD_BUILD_ID")) > 0
}

func onGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func onJenkins() bool {
	return len(os.Getenv("JENKINS_URL")) > 0
}

func onTeamCity() bool {
	return len(os.Getenv("TEAMCITY_VERSION")) > 0
}

func onTravisCI() bool {
	return os.Getenv("TRAVIS") == "true"
}

func onXcodeCloud() bool {
	return len(os.Getenv("CI_BUILD_ID")) > 0
}
