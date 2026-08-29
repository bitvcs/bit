package domain

import (
	"fmt"
	"net/url"
	"strings"
)

type NipaUrl struct {
	Url     string
	Host    string
	Org     string
	Project string
}

func ParseNipaUrl(urlString string) (*NipaUrl, error) {
	allowedSchemes := map[string]struct{}{
		"http":  {},
		"https": {},
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, NewUserError(fmt.Sprintf("url can not be parsed: %v", err))
	}

	if _, ok := allowedSchemes[parsedURL.Scheme]; !ok {
		return nil, NewUserError("invalid URL scheme, expected http or https")
	}

	pathSegments := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	if len(pathSegments) < 1 || len(pathSegments) > 2 {
		return nil, NewUserError("invalid URL path, expected /org/project")
	}

	org := pathSegments[0]
	project := pathSegments[0]
	if len(pathSegments) == 1 {
		org = "default"
	} else {
		project = pathSegments[1]
	}

	nipaUrl := &NipaUrl{
		Url:     urlString,
		Host:    parsedURL.Host,
		Org:     org,
		Project: project,
	}
	return nipaUrl, nil
}
