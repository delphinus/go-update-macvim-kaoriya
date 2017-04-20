package gumk

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"

	"github.com/Songmu/prompter"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gopkg.in/cheggaaa/pb.v1"
)

// Run will execute the process
func (g *Gumk) Run() error {
	eg, ctx := errgroup.WithContext(g.context)
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

	var p *pb.Pool
	eg.Go(func() error {
		var err error
		p, err = pb.StartPool(
			(<-dmgBar).Prefix("dmg    "),
			(<-appcastBar).Prefix("appcast"),
			(<-releaseBar).Prefix("release"),
		)
		if err != nil {
			return errors.Wrap(err, "error in StartPool")
		}
		return nil
	})

	err := eg.Wait()
	if p != nil {
		_ = p.Stop()
	}

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

	if err := f.save(); err != nil {
		return errors.Wrap(err, "error in save")
	}

	fmt.Println("saved successfully")

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
	fmt.Println("found:")
	fmt.Println(e)
	fmt.Println("to update:")
	fmt.Println(f.element)

	return prompter.YN("proceed?", false)
}
