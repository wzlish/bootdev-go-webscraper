package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	bufferSize := 5
	if len(os.Args) >= 3 {
		bufferSizeARG, err := strconv.Atoi(strings.Split(os.Args[2], ".")[0])
		if err != nil {
			fmt.Printf("invalid buffer size, expecting an int: %v", err)
			os.Exit(1)
		}
		bufferSize = bufferSizeARG
	}

	maxPages := 10
	if len(os.Args) >= 4 {
		maxPagesARG, err := strconv.Atoi(strings.Split(os.Args[3], ".")[0])
		if err != nil {
			fmt.Printf("invalid max pages, expecting an int: %v", err)
			os.Exit(1)
		}
		maxPages = maxPagesARG
	}

	baseUrl, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Errorf("invalid baseURL provided (%s): %v", os.Args[1], err)
		os.Exit(1)
	}

	crawler := crawler{
		seenPages:          make(map[string]int, 0),
		maxPages:           maxPages,
		baseURL:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, bufferSize),
		wg:                 &sync.WaitGroup{},
	}

	fmt.Println("Starting crawl of %s", crawler.baseURL)
	crawler.wg.Add(1)
	go crawler.crawlPage(baseUrl.String())
	crawler.wg.Wait()

	fmt.Println("\n--- Results ---")
	for url, count := range crawler.seenPages {
		fmt.Printf("%s (%d)\n", url, count)
	}

}
