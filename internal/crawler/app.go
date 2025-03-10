package crawler

import (
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

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	InternalRequest(link, baseURL string) (*http.Response, error)
	ExternalRequest(link string) (*http.Response, error)
}

type Crawler struct {
	httpClient HTTPClient
	baseUrl    string
}

func NewCrawler(httpClient HTTPClient) *Crawler {
	return &Crawler{
		httpClient: httpClient,
	}
}

func (c *Crawler) FindAllLinks(url string) {
	c.baseUrl = url
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
			c.processLink(l, linkCh, &activeTasksCount)
		}(link)
	}

	wg.Wait()
}

func (c *Crawler) processLink(currentLink string, linkCh chan string, numberOfActiveTask *int32) {
	response, err := c.makeRequest(currentLink)

	if err != nil {
		errorLog(currentLink, err)
	} else {
		c.handleResponse(currentLink, linkCh, response, numberOfActiveTask)
	}
}

func errorLog(selectedLink string, err error) {
	log.Printf("Link : %v | Error : %v ❌", selectedLink, err)
}

func (c *Crawler) handleResponse(selectedLink string, linkCh chan string, response *http.Response, numberOfActiveTask *int32) {
	if responseUtils.IsSuccess(response) {
		err := c.saveResponseLinks(response, linkCh, numberOfActiveTask)

		if err != nil {
			errorLog(selectedLink, err)
		}
	}

	logger.LogResponseResult(response)
}

func (c *Crawler) saveResponseLinks(response *http.Response, linkCh chan string, numberOfActiveTask *int32) error {
	if c.isInternal(response.Request.URL.String()) {
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

func (c *Crawler) makeRequest(link string) (*http.Response, error) {
	var err error
	var resp *http.Response

	if c.isInternal(link) {
		resp, err = c.httpClient.InternalRequest(link, c.baseUrl)
	} else {
		resp, err = c.httpClient.ExternalRequest(link)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Crawler) isInternal(link string) bool {
	return strings.HasPrefix(link, "/") || strings.HasPrefix(link, c.baseUrl)
}
