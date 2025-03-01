package main

import (
	"deadLinkParser/client"
	"deadLinkParser/data"
	"deadLinkParser/parsing"
	"log"
	"net/http"
	"strings"
	"sync"
)

var url string

func main() {
	url = "https://scrape-me.dreamsofcode.io/"

	response, err := makeRequest(url)
	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
	}

	linkCh := make(chan string, 100)
	var wg sync.WaitGroup

	err = saveResponseLinks(response, linkCh)

	if err != nil {
		log.Fatalf("Error while parsing html : %v", err)
	}

	for link := range linkCh {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			handleLink(l, linkCh)
		}(link)
	}

	wg.Wait()
	close(linkCh)
}

func handleLink(currentLink string, linkCh chan string) {
	
	response, err := makeRequest(currentLink)

	if err != nil {
		errorLog(currentLink, err)
	} else {
		handleResponse(currentLink, linkCh, response)
	}
}
func errorLog(selectedLink string, err error) {
	log.Printf("Link : %v | Error : %v ❌", selectedLink, err)
}

func handleResponse(selectedLink string, linkCh chan string, response *http.Response) {
	if response.StatusCode >= 200 && response.StatusCode < 400 {
		err := saveResponseLinks(response, linkCh)

		if err != nil {
			errorLog(selectedLink, err)
		}

		log.Printf("Link : %v | Status : %v ✅", selectedLink, response.Status)
	} else {
		log.Printf("Link : %v | Status : %v ❌", selectedLink, response.Status)
	}
}

func saveResponseLinks(response *http.Response, linkCh chan string) error {
	if isInternal(response.Request.URL.String()) {
		links, err := parsing.GetLinksFromResponse(response)

		if err != nil {
			return err
		}

		for _, link := range links {
			if data.CheckAndAddLink(link) {
				linkCh <- link
			}
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
