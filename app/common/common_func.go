package common

import (
	"regexp"
)

func GlobRoutePath(routePath string) (string, error) {
	regex, err := regexp.Compile(":[^/]+")
	if err != nil {
		return "", err
	}
	replacedPath := regex.ReplaceAllString(routePath, "*")
	return replacedPath, nil
}
