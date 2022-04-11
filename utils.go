package waldo

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Version() string {
	return fmt.Sprintf("%s %s (%s/%s)", agentName, agentVersion, detectPlatform(), detectArch())
}

//-----------------------------------------------------------------------------

func addIfNotEmpty(query *url.Values, key string, value string) {
	if len(key) > 0 && len(value) > 0 {
		query.Add(key, value)
	}
}

func appendIfNotEmpty(payload *string, key string, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}

	if len(*payload) > 0 {
		*payload += ","
	}

	*payload += fmt.Sprintf(`"%s":"%s"`, key, value)
}

func detectArch() string {
	arch := runtime.GOARCH

	switch arch {
	case "amd64":
		return "x86_64"

	default:
		return arch
	}
}

func detectPlatform() string {
	platform := runtime.GOOS

	switch platform {
	case "darwin":
		return "macOS"

	default:
		return strings.Title(platform)
	}
}

func determineBuildPayloadPath(workingPath, buildPath, buildSuffix string) string {
	buildName := filepath.Base(buildPath)

	switch buildSuffix {
	case "app":
		return filepath.Join(workingPath, buildName+".zip")

	default:
		return buildPath
	}
}

func determineWorkingPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("WaldoGoLib-%d", os.Getpid()))
}

func isDir(path string) bool {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return fi.Mode().IsDir()
}

func isRegular(path string) bool {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return fi.Mode().IsRegular()
}

func run(name string, args ...string) (string, string, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(name, args...)

	cmd.Stderr = &stderrBuffer
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	stderr := strings.TrimRight(stderrBuffer.String(), "\n")
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")

	return stdout, stderr, err
}

func validateBuildPath(buildPath string) (string, string, string, error) {
	if len(buildPath) == 0 {
		return "", "", "", errors.New("Empty build path")
	}

	buildPath, err := filepath.Abs(buildPath)

	if err != nil {
		return "", "", "", err
	}

	buildSuffix := strings.TrimPrefix(filepath.Ext(buildPath), ".")

	switch buildSuffix {
	case "apk":
		return buildPath, buildSuffix, "Android", nil

	case "app", "ipa":
		return buildPath, buildSuffix, "iOS", nil

	default:
		return "", "", "", fmt.Errorf("File extension of build at ‘%s’ is not recognized", buildPath)
	}
}

func validateUploadToken(uploadToken string) error {
	if len(uploadToken) == 0 {
		return errors.New("Empty upload token")
	}

	return nil
}

func zipDir(zipPath string, dirPath string, basePath string) error {
	err := os.Chdir(dirPath)

	if err != nil {
		return err
	}

	zipFile, err := os.Create(zipPath)

	if err != nil {
		return err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	walker := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		zipEntry, err := zipWriter.Create(path)

		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)

		return err
	}

	err = filepath.WalkDir(basePath, walker)

	err2 := zipWriter.Close()

	if err != nil {
		return err
	}

	return err2
}
