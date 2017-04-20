package gumk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type element struct {
	tag, version, dmg, appcast []byte
}

func (e element) String() string {
	return fmt.Sprintf(`  tag:     %s
  version: %s
  dmg:     %x
  appcast: %x`, string(e.tag), string(e.version), e.dmg, e.appcast)
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
	e.version = b[1]
	e.tag = b[2]
	e.dmg = convBytes(b[3])
	e.appcast = convBytes(a[2])
	return e, nil
}

func convBytes(b []byte) []byte {
	bb := make([]byte, 0, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		buf := bytes.NewBuffer(b[i : i+1])
		if i+1 < len(b) {
			_ = buf.WriteByte(b[i+1])
		}
		n, _ := strconv.ParseInt(buf.String(), 16, 16)
		bb = append(bb, byte(n))
	}
	return bb
}

var formulaReplace = `else
    version '%s:%s'
    sha256 '%x'`
var appcastReplace = `${1}%x`

func (f *Formula) save() error {
	formulaRepl := []byte(fmt.Sprintf(formulaReplace, string(f.version), string(f.tag), f.dmg))
	f.text = formulaRe.ReplaceAll(f.text, formulaRepl)
	appcastRepl := []byte(fmt.Sprintf(appcastReplace, f.appcast))
	f.text = appcastRe.ReplaceAll(f.text, appcastRepl)
	if err := ioutil.WriteFile(f.path, f.text, 0644); err != nil {
		return errors.Wrap(err, "error in WriteFile")
	}
	return nil
}
