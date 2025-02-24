package parsing

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func GetLinksFromResponse(response *http.Response) []string {
	htmlPage, err := html.Parse(response.Body)

	if err != nil {
		log.Fatalf("Error while parsing html : %v", err)
	}

	return getLinksFromHtmlPage(htmlPage)
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
