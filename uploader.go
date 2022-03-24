package waldo

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Uploader struct {
	userBuildPath   string
	userGitBranch   string
	userGitCommit   string
	userOverrides   map[string]string
	userUploadToken string
	userVariantName string
	userVerbose     bool

	arch             string
	buildPath        string
	buildPayloadPath string
	buildSuffix      string
	ciInfo           *CIInfo
	flavor           string
	gitInfo          *GitInfo
	platform         string
	validated        bool
	workingPath      string
}

//-----------------------------------------------------------------------------

func NewUploader(buildPath, uploadToken, variantName, gitCommit, gitBranch string, verbose bool, overrides map[string]string) *Uploader {
	return &Uploader{
		userBuildPath:   buildPath,
		userGitBranch:   gitBranch,
		userGitCommit:   gitCommit,
		userOverrides:   overrides,
		userUploadToken: uploadToken,
		userVariantName: variantName,
		userVerbose:     verbose}
}

//-----------------------------------------------------------------------------

func (u *Uploader) BuildPath() string {
	if u.validated {
		return u.buildPath
	}

	return u.userBuildPath
}

func (u *Uploader) BuildPayloadPath() string {
	return u.buildPayloadPath
}

func (u *Uploader) CIGitBranch() string {
	return u.ciInfo.GitBranch()
}

func (u *Uploader) CIGitCommit() string {
	return u.ciInfo.GitCommit()
}

func (u *Uploader) CIProvider() string {
	return u.ciInfo.Provider().String()
}

func (u *Uploader) GitAccess() string {
	return u.gitInfo.Access().String()
}

func (u *Uploader) GitBranch() string {
	return u.userGitBranch
}

func (u *Uploader) GitCommit() string {
	return u.userGitCommit
}

func (u *Uploader) InferredGitBranch() string {
	return u.gitInfo.Branch()
}

func (u *Uploader) InferredGitCommit() string {
	return u.gitInfo.Commit()
}

func (u *Uploader) UploadToken() string {
	return u.userUploadToken
}

func (u *Uploader) VariantName() string {
	return u.userVariantName
}

func (u *Uploader) Version() string {
	if !u.validated {
		return Version()
	}

	return fmt.Sprintf("%s %s (%s/%s)\n", agentName, agentVersion, u.platform, u.arch)
}

//-----------------------------------------------------------------------------

func (u *Uploader) Upload() error {
	err := os.RemoveAll(u.workingPath)

	if err == nil {
		err = os.MkdirAll(u.workingPath, 0755)
	}

	defer os.RemoveAll(u.workingPath)

	if err == nil {
		err = u.createBuildPayload()
	}

	if err == nil {
		err = u.uploadBuild()
	}

	if err != nil {
		u.uploadError(err)
	}

	return err
}

func (u *Uploader) Validate() error {
	if u.validated {
		return nil
	}

	err := validateUploadToken(u.userUploadToken)

	if err != nil {
		return err
	}

	buildPath, buildSuffix, flavor, err := validateBuildPath(u.userBuildPath)

	if err != nil {
		// u.uploadError(err) ???

		return err
	}

	workingPath := determineWorkingPath()

	u.arch = detectArch()
	u.buildPath = buildPath
	u.buildPayloadPath = determineBuildPayloadPath(workingPath, buildPath, buildSuffix)
	u.buildSuffix = buildSuffix
	u.ciInfo = DetectCIInfo(true)
	u.flavor = flavor
	u.gitInfo = InferGitInfo(u.ciInfo.SkipCount())
	u.platform = detectPlatform()
	u.validated = true
	u.workingPath = workingPath

	return nil
}

//-----------------------------------------------------------------------------

func (u *Uploader) authorization() string {
	return fmt.Sprintf("Upload-Token %s", u.userUploadToken)
}

func (u *Uploader) buildContentType() string {
	switch u.buildSuffix {
	case "app":
		return "application/zip"

	default:
		return "application/octet-stream"
	}
}

func (u *Uploader) checkBuildStatus(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	bodyString := string(body)

	statusRegex := regexp.MustCompile(`"status":([0-9]+)`)
	statusMatches := statusRegex.FindStringSubmatch(bodyString)

	if len(statusMatches) > 0 { // status is numeric _only_ on failure
		var status = 0

		status, err = strconv.Atoi(statusMatches[1])

		if err != nil {
			return err
		}

		if status == 401 {
			return fmt.Errorf("Upload token is invalid or missing!")
		}

		if status < 200 || status > 299 {
			return fmt.Errorf("Unable to upload build to Waldo, HTTP status: %d", status)
		}
	}

	return nil
}

func (u *Uploader) createBuildPayload() error {
	parentPath := filepath.Dir(u.buildPath)
	buildName := filepath.Base(u.buildPath)

	switch u.buildSuffix {
	case "app":
		if !isDir(u.buildPath) {
			return fmt.Errorf("Unable to read build at ‘%s’", u.buildPath)
		}

		return zipDir(u.buildPayloadPath, parentPath, buildName)

	default:
		if !isRegular(u.buildPath) {
			return fmt.Errorf("Unable to read build at ‘%s’", u.buildPath)
		}

		return nil
	}
}

func (u *Uploader) dumpRequest(req *http.Request, body bool) {
	if u.userVerbose {
		dump, err := httputil.DumpRequestOut(req, body)

		if err == nil {
			fmt.Printf("\n--- Request ---\n%s\n", dump)
		}
	}
}

func (u *Uploader) dumpResponse(resp *http.Response, body bool) {
	if u.userVerbose {
		dump, err := httputil.DumpResponse(resp, body)

		if err == nil {
			fmt.Printf("\n--- Response ---\n%s\n", dump)
		}
	}
}

func (u *Uploader) errorContentType() string {
	return "application/json"
}

func (u *Uploader) makeBuildURL() string {
	buildURL := u.userOverrides["apiBuildEndpoint"]

	if len(buildURL) == 0 {
		buildURL = defaultAPIBuildEndpoint
	}

	query := make(url.Values)

	addIfNotEmpty(&query, "agentName", agentName)
	addIfNotEmpty(&query, "agentVersion", agentVersion)
	addIfNotEmpty(&query, "arch", u.arch)
	addIfNotEmpty(&query, "ci", u.ciInfo.Provider().String())
	addIfNotEmpty(&query, "ciGitBranch", u.ciInfo.GitBranch())
	addIfNotEmpty(&query, "ciGitCommit", u.ciInfo.GitCommit())
	addIfNotEmpty(&query, "flavor", u.flavor)
	addIfNotEmpty(&query, "gitAccess", u.gitInfo.Access().String())
	addIfNotEmpty(&query, "gitBranch", u.gitInfo.Branch())
	addIfNotEmpty(&query, "gitCommit", u.gitInfo.Commit())
	addIfNotEmpty(&query, "platform", u.platform)
	addIfNotEmpty(&query, "userGitBranch", u.userGitBranch)
	addIfNotEmpty(&query, "userGitCommit", u.userGitCommit)
	addIfNotEmpty(&query, "variantName", u.userVariantName)
	addIfNotEmpty(&query, "wrapperName", u.userOverrides["wrapperName"])
	addIfNotEmpty(&query, "wrapperVersion", u.userOverrides["wrapperVersion"])

	buildURL += "?" + query.Encode()

	return buildURL
}

func (u *Uploader) makeErrorPayload(err error) string {
	payload := ""

	appendIfNotEmpty(&payload, "agentName", agentName)
	appendIfNotEmpty(&payload, "agentVersion", agentVersion)
	appendIfNotEmpty(&payload, "arch", u.arch)
	appendIfNotEmpty(&payload, "ci", u.ciInfo.Provider().String())
	appendIfNotEmpty(&payload, "ciGitBranch", u.ciInfo.GitBranch())
	appendIfNotEmpty(&payload, "ciGitCommit", u.ciInfo.GitCommit())
	appendIfNotEmpty(&payload, "message", err.Error())
	appendIfNotEmpty(&payload, "platform", u.platform)
	appendIfNotEmpty(&payload, "wrapperName", u.userOverrides["wrapperName"])
	appendIfNotEmpty(&payload, "wrapperVersion", u.userOverrides["wrapperVersion"])

	payload = "{" + payload + "}"

	return payload
}

func (u *Uploader) makeErrorURL() string {
	errorURL := u.userOverrides["apiErrorEndpoint"]

	if len(errorURL) == 0 {
		errorURL = defaultAPIErrorEndpoint
	}

	return errorURL
}

func (u *Uploader) uploadBuild() error {
	url := u.makeBuildURL()

	file, err := os.Open(u.buildPayloadPath)

	if err != nil {
		return fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url)
	}

	defer file.Close()

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, file)

	if err != nil {
		return fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url)
	}

	req.Header.Add("Authorization", u.authorization())
	req.Header.Add("Content-Type", u.buildContentType())
	req.Header.Add("User-Agent", u.userAgent())

	u.dumpRequest(req, false)

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url)
	}

	u.dumpResponse(resp, true)

	defer resp.Body.Close()

	return u.checkBuildStatus(resp)
}

func (u *Uploader) uploadError(err error) error {
	url := u.makeErrorURL()
	body := u.makeErrorPayload(err)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", u.authorization())
	req.Header.Add("Content-Type", u.errorContentType())
	req.Header.Add("User-Agent", u.userAgent())

	// u.dumpRequest(req, true)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// u.dumpResponse(resp, true)

	return nil
}

func (u *Uploader) userAgent() string {
	ci := u.ciInfo.Provider().String()

	if ci == "Unknown" {
		ci = "Go CLI" // hack for now…
	}

	version := u.userOverrides["wrapperVersion"]

	if len(version) == 0 {
		version = agentVersion
	}

	return fmt.Sprintf("Waldo %s/%s v%s", ci, u.flavor, version)
}
