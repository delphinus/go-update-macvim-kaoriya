package gumk

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/pkg/errors"
)

type element struct {
	tag, version, dmg, appcast []byte
}

// Formula is a struct for formulas
type Formula struct {
	path string
	text []byte
	element
}

// NewFormula returns a new Formula instance
func NewFormula(path, tag string, version, dmg, appcast []byte) Formula {
	return Formula{
		path: path,
		element: element{
			tag:     []byte(tag),
			version: version,
			dmg:     dmg,
			appcast: appcast,
		},
	}
}

var formulaRe = regexp.MustCompile(`else
    version '(\d\.\d):(\d+)'
    sha256 '([\da-f]+)'`)
var appcastRe = regexp.MustCompile(`(checkpoint: ')([\da-f]+)`)

func (f *Formula) read() (element, error) {
	e := element{}
	text, err := ioutil.ReadFile(f.path)
	if err != nil {
		return e, errors.Wrap(err, "error in ReadFile")
	}

	b := formulaRe.FindSubmatch(text)
	if len(b) != 4 {
		return e, errors.Wrap(err, "cannot find elements")
	}

	a := appcastRe.FindSubmatch(text)
	if len(a) != 3 {
		return e, errors.Wrap(err, "cannot find appcast")
	}

	f.text = text
	e.appcast = a[2]
	e.version = b[1]
	e.tag = b[2]
	e.dmg = b[3]
	return e, nil
}

var formulaReplace = `else
    version '%s:%s'
    sha256 '%s'`
var appcastReplace = `$1%s`

func (f *Formula) save(e element) error {
	formulaRepl := []byte(fmt.Sprintf(formulaReplace, e.version, e.tag, e.dmg))
	f.text = formulaRe.ReplaceAll(f.text, formulaRepl)
	appcastRepl := []byte(fmt.Sprintf(appcastReplace, e.appcast))
	f.text = appcastRe.ReplaceAll(f.text, appcastRepl)
	if err := ioutil.WriteFile(f.path, f.text, 0644); err != nil {
		return errors.Wrap(err, "error in WriteFile")
	}
	return nil
}
