package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/delphinus/go-update-macvim-kaoriya"
)

func main() {
	flag.Parse()
	tag := flag.Arg(0)
	if tag == "" {
		fmt.Fprintln(os.Stderr, "tag needed to execute")
		os.Exit(1)
	}
	app, err := gumk.New(
		gumk.WithContext(context.Background()),
		gumk.WithTag(tag),
		gumk.WithHTTPClient(http.DefaultClient),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error ocurred: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error ocurred: %v\n", err)
		os.Exit(1)
	}
}
