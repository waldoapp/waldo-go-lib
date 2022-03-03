package waldo

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
)

type Triggerer struct {
	userOverrides   map[string]string
	userRuleName    string
	userUploadToken string
	userVerbose     bool

	arch      string
	ci        string
	platform  string
	validated bool
}

//-----------------------------------------------------------------------------

func NewTriggerer(uploadToken, ruleName string, verbose bool, overrides map[string]string) *Triggerer {
	return &Triggerer{
		userOverrides:   overrides,
		userRuleName:    ruleName,
		userUploadToken: uploadToken,
		userVerbose:     verbose}
}

//-----------------------------------------------------------------------------

func (t *Triggerer) RuleName() string {
	return t.userRuleName
}

func (t *Triggerer) UploadToken() string {
	return t.userUploadToken
}

func (t *Triggerer) Version() string {
	if !t.validated {
		return Version()
	}

	return fmt.Sprintf("%s %s (%s/%s)\n", agentName, agentVersion, t.platform, t.arch)
}

//-----------------------------------------------------------------------------

func (t *Triggerer) Perform() error {
	return t.triggerRun()
}

func (t *Triggerer) Validate() error {
	if t.validated {
		return nil
	}

	err := validateUploadToken(t.userUploadToken)

	if err != nil {
		return err
	}

	t.arch = detectArch()
	t.ci = detectCI()
	t.platform = detectPlatform()
	t.validated = true

	return nil
}

//-----------------------------------------------------------------------------

func (t *Triggerer) authorization() string {
	return fmt.Sprintf("Upload-Token %s", t.userUploadToken)
}

func (t *Triggerer) checkTriggerStatus(resp *http.Response) error {
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
			return fmt.Errorf("Unable to trigger run on Waldo, HTTP status: %d", status)
		}
	}

	return nil
}

func (t *Triggerer) contentType() string {
	return "application/json"
}

func (t *Triggerer) dumpRequest(req *http.Request, body bool) {
	if t.userVerbose {
		dump, err := httputil.DumpRequestOut(req, body)

		if err == nil {
			fmt.Printf("\n--- Request ---\n%s\n", dump)
		}
	}
}

func (t *Triggerer) dumpResponse(resp *http.Response, body bool) {
	if t.userVerbose {
		dump, err := httputil.DumpResponse(resp, body)

		if err == nil {
			fmt.Printf("\n--- Response ---\n%s\n", dump)
		}
	}
}

func (t *Triggerer) makePayload() string {
	payload := ""

	appendIfNotEmpty(&payload, "agentName", agentName)
	appendIfNotEmpty(&payload, "agentVersion", agentVersion)
	appendIfNotEmpty(&payload, "arch", t.arch)
	appendIfNotEmpty(&payload, "ci", t.ci)
	appendIfNotEmpty(&payload, "platform", t.platform)
	appendIfNotEmpty(&payload, "ruleName", t.userRuleName)
	appendIfNotEmpty(&payload, "wrapperName", t.userOverrides["wrapperName"])
	appendIfNotEmpty(&payload, "wrapperVersion", t.userOverrides["wrapperVersion"])

	payload = "{" + payload + "}"

	return payload
}

func (t *Triggerer) makeURL() string {
	triggerURL := t.userOverrides["apiTriggerEndpoint"]

	if len(triggerURL) == 0 {
		triggerURL = defaultAPITriggerEndpoint
	}

	return triggerURL
}

func (t *Triggerer) triggerRun() error {
	url := t.makeURL()
	body := t.makePayload()

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return fmt.Errorf("Unable to trigger run on Waldo, error: %v, url: %s", err, url)
	}

	req.Header.Add("Authorization", t.authorization())
	req.Header.Add("Content-Type", t.contentType())
	req.Header.Add("User-Agent", t.userAgent())

	t.dumpRequest(req, true)

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Unable to trigger run on Waldo, error: %v, url: %s", err, url)
	}

	t.dumpResponse(resp, true)

	defer resp.Body.Close()

	return t.checkTriggerStatus(resp)
}

func (t *Triggerer) userAgent() string {
	ci := t.ci

	if len(ci) == 0 {
		ci = "Go CLI" // hack for now…
	}

	version := t.userOverrides["wrapperVersion"]

	if len(version) == 0 {
		version = agentVersion
	}

	return fmt.Sprintf("Waldo %s v%s", ci, version)
}
