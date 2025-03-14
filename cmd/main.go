package main

import (
	"context"
	"deadLinkParser/internal/crawler"
	"deadLinkParser/internal/http/client"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "deadLinkParser",
		Usage: "call all links in your website",
		Action: func(ctx context.Context, command *cli.Command) error {
			httpClient := client.NewRealHTTPClient()
			crawler := crawler.NewCrawler(httpClient)
			crawler.FindAllLinks(command.Args().Get(0))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
