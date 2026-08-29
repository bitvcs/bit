package main

import (
	"fmt"
	"net/url"
	"strings"
)

type nipaUrl struct {
	Url     string
	Host    string
	Org     string
	Project string
}

func parseURL(urlString string) (*nipaUrl, error) {
	allowedSchemes := map[string]struct{}{
		"http":  {},
		"https": {},
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if _, ok := allowedSchemes[parsedURL.Scheme]; !ok {
		return nil, fmt.Errorf("invalid URL scheme, expected http or https")
	}

	pathSegments := strings.Split(parsedURL.Path[1:], "/")
	if len(pathSegments) != 2 {
		return nil, fmt.Errorf("invalid URL path, expected /org/project")
	}
	fmt.Println(pathSegments)

	nipaUrl := &nipaUrl{
		Url:  urlString,
		Host: parsedURL.Host,
	}
	return nipaUrl, nil
}
