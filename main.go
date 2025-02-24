package main

import (
	"deadLinkParser/data"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

func main() {
	url := "https://scrape-me.dreamsofcode.io/"

	appData := data.InitialiseAppData(url)

	links := getHrefLinks(url)

	for _, link := range links {
		data.AddLinkFound(link, &appData)
	}

	index := 0

	for index < len(appData.FoundLinks) {
		if !appData.HasUncheckedLink() {
			break
		}

		selectedLink := appData.FoundLinks[index].Link
		if strings.HasPrefix(selectedLink, "/") {
			fullUrl := buildUrl(selectedLink, url)

			newLinks := getHrefLinks(fullUrl)

			for _, link := range newLinks {
				data.AddLinkFound(link, &appData)
			}
		}

		index++

		log.Print(appData)
	}

}

func getHrefLinks(url string) []string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error requesting url '%v' : %v", url, err)
	}

	content, err := html.Parse(resp.Body)

	var links []string

	for element := range content.Descendants() {
		if element.Type == html.ElementNode && element.Data == "a" {
			for _, attribute := range element.Attr {
				if attribute.Key == "href" {
					links = append(links, attribute.Val)
				}
			}
		}
	}

	return links
}

func buildUrl(link string, baseUrl string) string {
	return baseUrl + link[1:]
}
