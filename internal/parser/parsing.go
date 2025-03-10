package parser

import (
	"net/http"

	"golang.org/x/net/html"
)

func GetLinksFromResponse(response *http.Response) ([]string, error) {
	htmlPage, err := html.Parse(response.Body)

	if err != nil {
		return nil, err
	}

	return getLinksFromHtmlPage(htmlPage), nil
}

func getLinksFromHtmlPage(htmlPage *html.Node) []string {
	var links []string

	for element := range htmlPage.Descendants() {
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
