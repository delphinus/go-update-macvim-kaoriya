package gumk

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Songmu/prompter"
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

	version, err := g.findVersion(release.Bytes())
	if err != nil {
		return errors.Wrap(err, "error in findVersion")
	}

	f := NewFormula(g.formula, g.tag, version, dmg.Sum(nil), appcast.Sum(nil))
	e, err := f.read()
	if err != nil {
		return errors.Wrap(err, "error in read")
	}

	if !g.confirmProceed(f, e) {
		return errors.New("cancelled")
	}

	if err := f.save(e); err != nil {
		return errors.Wrap(err, "error in save")
	}

	fmt.Print("saved successfully")

	return nil
}

var versionRe = regexp.MustCompile(`Vim (\d\.\d)`)

func (g *Gumk) findVersion(s []byte) ([]byte, error) {
	sm := versionRe.FindSubmatch(s)
	if len(sm) < 2 || len(sm[1]) != 3 {
		return nil, errors.New("cannot find version string")
	}
	return sm[1], nil
}

func (g *Gumk) confirmProceed(f Formula, e element) bool {
	fmt.Printf(`found:
  tag:     %s
  version: %s
  dmg:     %s
  appcast: %s

`, string(e.tag), string(e.version), string(e.dmg), string(e.appcast))
	fmt.Printf(`to update:
  tag:     %s
  version: %s
  dmg:     %x
  appcast: %x

`, f.tag, string(f.version), f.dmg, f.appcast)

	return prompter.YN("proceed?", false)
}
