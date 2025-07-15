package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}

	req.Header.Set("User-Agent", "bootdev-bootie-crawler/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getting response failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("non 2XX http status code (%d)", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to ReadAll body: %v", err)

	}

	return string(respBody), nil

}
