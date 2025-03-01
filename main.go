package main

import (
	"context"
	"deadLinkParser/app"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	cmd := &cli.Command{
		Name:  "deadLinkParser",
		Usage: "call all links in your website",
		Action: func(ctx context.Context, command *cli.Command) error {
			app.FindAllLinks(command.Args().Get(0))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
