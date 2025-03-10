package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	client, err := dagger.Connect(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to Dagger: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	src := client.Host().Directory(".")

	golang := client.Container().
		From("golang:1.24").
		WithDirectory("/src", src).
		WithWorkdir("/src")

	golang = golang.WithExec([]string{"go", "mod", "download"})

	out, err := golang.WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Test output:")
	fmt.Println(out)
}
