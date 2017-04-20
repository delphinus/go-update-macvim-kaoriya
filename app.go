package gumk

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// Gumk is a struct of the app
type Gumk struct {
	context               context.Context
	tag, formula, appcast string
	dmg, release          func() string
	httpClient            *http.Client
}

// New returns an instance of the app
func New(opts ...Option) (*Gumk, error) {
	g := &Gumk{}

	opts = append(opts,
		WithFormula(),
		WithAppcast(),
	)

	for _, o := range opts {
		o.apply(g)
	}

	if g.context == nil {
		return nil, errors.New("cannot detect context")
	}

	if g.tag == "" {
		return nil, errors.New("cannot detect tag")
	}

	if g.httpClient == nil {
		return nil, errors.New("cannot detect HTTP client")
	}

	WithDMG(g.tag).apply(g)
	WithRelease(g.tag).apply(g)

	return g, nil
}
