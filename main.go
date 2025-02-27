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

	response, err := makeRequest(url)
	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
	}

	err = saveResponseLinks(response, &appData)

	if err != nil {
		log.Fatalf("Error while parsing html : %v", err)
	}

	index := 0

	for index < len(appData.FoundLinks) {
		selectedLink := appData.FoundLinks[index].Link

		response, err = makeRequest(selectedLink)

		if err != nil {
			errorLog(selectedLink, err)
		} else {
			handleResponse(selectedLink, &appData, response)
		}

		index++
	}

}

func errorLog(selectedLink string, err error) {
	log.Printf("Link : %v | Error : %v ❌", selectedLink, err)
}

func handleResponse(selectedLink string, appData *data.AppData, response *http.Response) {
	if response.StatusCode >= 200 && response.StatusCode < 400 {
		data.UpdateLink(selectedLink, appData, true, true)

		err := saveResponseLinks(response, appData)

		if err != nil {
			errorLog(selectedLink, err)
		}

		log.Printf("Link : %v | Status : %v ✅", selectedLink, response.Status)
	} else {
		log.Printf("Link : %v | Status : %v ❌", selectedLink, response.Status)
	}
}

func saveResponseLinks(response *http.Response, appData *data.AppData) error {
	if isInternal(response.Request.URL.String()) {
		links, err := parsing.GetLinksFromResponse(response)

		if err != nil {
			return err
		}

		for _, link := range links {
			data.AddLinkFound(link, appData)
		}

	}
	return nil
}

func makeRequest(link string) (*http.Response, error) {
	var err error
	var resp *http.Response

	if isInternal(link) {
		resp, err = client.InternalRequest(link, url)
	} else {
		resp, err = client.ExternalRequest(link)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func isInternal(link string) bool {
	return strings.HasPrefix(link, "/") || strings.HasPrefix(link, url)
}
