package crawler

import (
	"deadLinkParser/internal/http/client"
	responseUtils "deadLinkParser/internal/http/utils"
	"deadLinkParser/internal/logger"
	parsing "deadLinkParser/internal/parser"
	data "deadLinkParser/internal/storage"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var baseUrl string

func FindAllLinks(url string) {

	baseUrl = url
	linkCh := make(chan string, 100)
	var wg sync.WaitGroup

	var activeTasksCount int32 = 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(linkCh)

		for {
			time.Sleep(100 * time.Millisecond)
			if atomic.LoadInt32(&activeTasksCount) == 0 {
				break
			}
		}
	}()

	atomic.AddInt32(&activeTasksCount, 1)
	linkCh <- url

	for link := range linkCh {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			defer atomic.AddInt32(&activeTasksCount, -1)
			processLink(l, linkCh, &activeTasksCount)
		}(link)
	}

	wg.Wait()
}
func processLink(currentLink string, linkCh chan string, numberOfActiveTask *int32) {

	response, err := makeRequest(currentLink)

	if err != nil {
		errorLog(currentLink, err)
	} else {
		handleResponse(currentLink, linkCh, response, numberOfActiveTask)
	}
}
func errorLog(selectedLink string, err error) {
	log.Printf("Link : %v | Error : %v ❌", selectedLink, err)
}

func handleResponse(selectedLink string, linkCh chan string, response *http.Response, numberOfActiveTask *int32) {
	if responseUtils.IsSuccess(response) {
		err := saveResponseLinks(response, linkCh, numberOfActiveTask)

		if err != nil {
			errorLog(selectedLink, err)
		}
	}

	logger.LogResponseResult(response)
}

func saveResponseLinks(response *http.Response, linkCh chan string, numberOfActiveTask *int32) error {
	if isInternal(response.Request.URL.String()) {
		links, err := parsing.GetLinksFromResponse(response)

		if err != nil {
			return err
		}

		for _, link := range links {
			if data.CheckAndAddLink(link) {
				atomic.AddInt32(numberOfActiveTask, 1)
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
		resp, err = client.InternalRequest(link, baseUrl)
	} else {
		resp, err = client.ExternalRequest(link)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func isInternal(link string) bool {
	return strings.HasPrefix(link, "/") || strings.HasPrefix(link, baseUrl)
}
