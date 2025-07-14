package main

import (
	"errors"
	"net/url"
)

var ErrInvalidScheme = errors.New("Not a http/https url")
var ErrInvalidHost = errors.New("Input url does not contain a valid host")

func normalizeURL(inputURL string) (string, error) {

	url, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !(url.Scheme == "http" || url.Scheme == "https") {
		return "", ErrInvalidScheme
	}

	if len(url.Scheme) == 0 {
		return "", ErrInvalidHost
	}

	output := url.Host
	if len(url.Path) > 0 {
		output += url.Path
	}
	if len(url.RawQuery) > 0 {
		output += "?" + url.RawQuery
	}
	return output, nil

}
