package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	gotHTML, err := getHTML(baseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(gotHTML)

}
