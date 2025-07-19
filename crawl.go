package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
)

type crawler struct {
	seenPages          map[string]int
	maxPages           int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (c *crawler) addPageVisit(normalizedURL string) (isFirst bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.seenPages[normalizedURL]
	if ok {
		c.seenPages[normalizedURL]++
		return false
	}
	c.seenPages[normalizedURL] = 1
	return true
}

func (c *crawler) numSeenPages() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.seenPages)
}

func (c *crawler) crawlPage(rawCurrentURL string) {

	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
		c.wg.Done()
	}()

	if c.numSeenPages() >= c.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("unable to parse current url (%s): %v", rawCurrentURL, err)
		return
	}

	if currentURL.Host != c.baseURL.Host {
		return
	}

	normalizedCurrent, err := normalizeURL(currentURL.String())
	if err != nil {
		log.Printf("unable to normalize current url (%s): %v", currentURL.String(), err)
		return
	}

	firstPage := c.addPageVisit(normalizedCurrent)
	if !firstPage {
		return
	}
	fmt.Printf("Crawling /%s\n", currentURL.Path)

	currentHTML, err := getHTML(currentURL.String())
	if err != nil {
		log.Printf("unable to get html for %s: %v", currentURL.String(), err)
		return
	}

	foundURLS, err := getURLsFromHTML(currentHTML, c.baseURL.String())
	if err != nil {
		log.Printf("unable to get urls for %s: %v", currentURL.String(), err)
		return
	}

	for _, url := range foundURLS {
		c.wg.Add(1)
		go c.crawlPage(url)
	}

}
