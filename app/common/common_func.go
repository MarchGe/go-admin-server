package common

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
)

func GlobRoutePath(routePath string) (string, error) {
	regex, err := regexp.Compile(":[^/]+")
	if err != nil {
		return "", err
	}
	replacedPath := regex.ReplaceAllString(routePath, "*")
	return replacedPath, nil
}

func GetBash() (string, error) {
	unsupportedPlatforms := []string{"windows"} // pty not support windows platform
	if slices.Contains(unsupportedPlatforms, runtime.GOOS) {
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	if _, err := exec.LookPath("bash"); err != nil {
		return "", fmt.Errorf("bash command not found, %w", err)
	}
	return "bash", nil
}
