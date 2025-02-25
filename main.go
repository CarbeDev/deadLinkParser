package main

import (
	"deadLinkParser/client"
	"deadLinkParser/data"
	"deadLinkParser/parsing"
	"log"
	"net/http"
	"strings"
)

var url string

func main() {
	url = "https://scrape-me.dreamsofcode.io/"
	appData := data.InitialiseAppData(url)

	response := makeRequest(url)
	log.Printf("%v", response)
	links := parsing.GetLinksFromResponse(response)
	saveLinks(links, &appData)

	index := 0

	for index < len(appData.FoundLinks) {
		if !appData.HasUncheckedLink() {
			break
		}

		selectedLink := appData.FoundLinks[index].Link
		if isInternal(selectedLink) {
			response = makeRequest(selectedLink)
			handleRequest(selectedLink, appData, response)
		}

		index++

		log.Print(appData)
	}

}

func handleRequest(selectedLink string, appData data.AppData, response *http.Response) {
	if response.StatusCode >= 200 && response.StatusCode < 400 {
		data.UpdateLink(selectedLink, &appData, true, true)

		links := parsing.GetLinksFromResponse(response)
		saveLinks(links, &appData)
	}
}

func makeRequest(link string) *http.Response {
	var err error
	var resp *http.Response

	if isInternal(link) {
		resp, err = client.InternalRequest(link, url)
	} else {
		resp, err = client.ExternalRequest(link)
	}

	if err != nil {
		log.Fatalf("Error requesting url '%v' : %v", url, err)
	}

	return resp
}

func isInternal(link string) bool {
	return strings.HasPrefix(link, "/")
}

func saveLinks(links []string, appData *data.AppData) {
	for _, link := range links {
		data.AddLinkFound(link, appData)
	}
}
