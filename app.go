package gumk

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gopkg.in/cheggaaa/pb.v1"
)

// Gumk is a struct of the app
type Gumk struct {
	context               context.Context
	tag, formula, appcast string
	dmg, release          func() string
	httpClient            *http.Client
}

// New returns an instance of the app
func New(options ...Option) *Gumk {
	g := Gumk{}

	for _, o := range options {
		o.apply(&g)
	}

	return &g
}

// Run will execute the process
func (g *Gumk) Run() error {
	eg, ctx := errgroup.WithContext(context.Background())
	g.context = ctx
	dmg := sha256.New()
	appcast := sha256.New()
	release := bytes.NewBuffer(nil)

	dmgBar := make(chan *pb.ProgressBar)
	appcastBar := make(chan *pb.ProgressBar)
	releaseBar := make(chan *pb.ProgressBar)

	eg.Go(func() error { return g.fetch(g.dmg(), dmg, dmgBar) })
	eg.Go(func() error { return g.fetch(g.appcast, appcast, appcastBar) })
	eg.Go(func() error { return g.fetch(g.release(), release, releaseBar) })

	p, err := pb.StartPool(
		(<-dmgBar).Prefix("dmg    "),
		(<-appcastBar).Prefix("appcast"),
		(<-releaseBar).Prefix("release"),
	)
	if err != nil {
		return errors.Wrap(err, "error in StartPool") // do not die here
	}

	err = eg.Wait()
	_ = p.Stop()

	if err != nil {
		return errors.Wrap(err, "error in Wait")
	}

	fmt.Printf("dmgHash: %x\n", dmg.Sum(nil))
	fmt.Printf("appcastHash: %x\n", appcast.Sum(nil))
	fmt.Printf("release: %d\n", len(release.String()))
	return nil
}
