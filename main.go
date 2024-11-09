package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/url"
	"os"
	"strings"
)

func extractLinks(htmlString string) (map[string]int, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, err
	}

	res := make(map[string]int)
	
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if _, err := url.Parse(attr.Val); err == nil {
						res[attr.Val]++
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return res, nil
}

func main() {
	flag.Parse()
	htmlStrings := flag.Args()

	if len(htmlStrings) == 0 {
		// Read from stdin if no arguments are provided
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Read all input from stdin
			stdinBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				return
			}
			htmlStrings = append(htmlStrings, string(stdinBytes))
		} else {
			fmt.Println("No input provided. Use -h for help.")
			return
		}
	}

	for _, htmlString := range htmlStrings {
		m, err := extractLinks(htmlString)
		if err != nil {
			fmt.Println("Error extracting links:", err)
			return
		}

		fmt.Println("Link - Count:")
		for link, count := range m {
			fmt.Printf("%s - %d\n", link, count)
		}
	}

}
