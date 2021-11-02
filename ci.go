package waldo

import (
	"os"
)

func detectCI() string {
	switch {
	case onAppCenter():
		return "App Center"

	case onAzureDevOps():
		return "Azure DevOps"

	case onBitrise():
		return "Bitrise"

	case onBuddybuild():
		return "buddybuild"

	case onCircleCI():
		return "CircleCI"

	case onCodeBuild():
		return "CodeBuild"

	case onGitHubActions():
		return "GitHub Actions"

	case onJenkins():
		return "Jenkins"

	case onTeamCity():
		return "TeamCity"

	case onTravisCI():
		return "Travis CI"

	case onXcodeCloud():
		return "Xcode Cloud"

	default:
		return ""
	}
}

func getSkipCount() int {
	if onGitHubActions() &&
		os.Getenv("GITHUB_EVENT_NAME") == "pull_request" {
		return 1
	}

	return 0
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

func onBuddybuild() bool {
	return len(os.Getenv("BUDDYBUILD_BUILD_ID")) > 0
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
