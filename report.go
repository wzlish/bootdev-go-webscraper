package main

import (
	"fmt"
	"sort"
)

type Page struct {
	url   string
	count int
}

func orderSeenPagesDesc(pages map[string]int) []Page {

	pagesDesc := make([]Page, 0, len(pages))
	for k, v := range pages {
		pagesDesc = append(pagesDesc, Page{k, v})
	}

	sort.Slice(pagesDesc, func(i, j int) bool {
		return pagesDesc[i].count > pagesDesc[j].count
	})

	return pagesDesc
}

func printReport(url string, report []Page) {
	fmt.Printf("\n=============================\nREPORT for %s\n=============================\n", url)
	for _, page := range report {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.url)
	}
}
