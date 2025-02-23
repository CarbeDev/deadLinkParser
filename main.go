package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func main() {
	url := "https://scrape-me.dreamsofcode.io/"

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error requesting url '%v' : %v", url, err)
	}

	//readable, _ := io.ReadAll(resp.Body)
	//log.Print(string(readable))

	content, err := html.Parse(resp.Body)

	for element := range content.Descendants() {
		if element.Type == html.ElementNode && element.Data == "a" {
			for _, attribute := range element.Attr {
				if attribute.Key == "href" {
					fmt.Println(attribute.Val)
				}
			}
		}
	}

}
