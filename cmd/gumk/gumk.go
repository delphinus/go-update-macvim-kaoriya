package main

import (
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
	app := gumk.New(
		gumk.WithTag(tag),
		gumk.WithFormula(),
		gumk.WithDMG(tag),
		gumk.WithRelease(tag),
		gumk.WithAppcast(),
		gumk.WithHTTPClient(http.DefaultClient),
	)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error ocurred: %v\n", err)
		os.Exit(1)
	}
}
