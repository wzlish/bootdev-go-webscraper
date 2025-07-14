package main

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var ErrInvalidBaseURL = errors.New("invalid or empty base URL provided")

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

	if len(rawBaseURL) == 0 {
		return nil, ErrInvalidBaseURL
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil || baseURL == nil || len(baseURL.String()) == 0 || !(baseURL.Scheme == "http" || baseURL.Scheme == "https") {
		return nil, ErrInvalidBaseURL
	}

	htmlReader := strings.NewReader(htmlBody)
	htmlNodes, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}

	var foundHrefs = make([]string, 0)

	for thisNode := range htmlNodes.Descendants() {

		if thisNode.Type == html.ElementNode && thisNode.Data == "a" {
			for _, attr := range thisNode.Attr {
				if attr.Key == "href" {

					thisURL, err := url.Parse(attr.Val)
					if err != nil {
						log.Printf("Error parsing href '%s': %v", attr.Val, err)
						continue
					}

					if thisURL.Scheme != "" && thisURL.Scheme != "http" && thisURL.Scheme != "https" {
						log.Printf("Skipping href with unsupported scheme '%s': %s", thisURL.Scheme, attr.Val)
						continue
					}

					foundHrefs = append(foundHrefs, baseURL.ResolveReference(thisURL).String())
					continue
				}
			}

		}
	}

	return foundHrefs, nil
}
