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

var re = regexp.MustCompile(`else
    version '(?P<version>\d\.\d):(?P<tag>\d+)'
    sha256 '(?P<dmg>[\da-f]+)'
(?P<skipped>(?s:.*))checkpoint: '(?P<appcast>[\da-f]+)`)

var tmpl = `else
    version '%s:%s'
    sha256 '%x'
${skipped}checkpoint: '%x`

func (f *Formula) findMatches() (map[string][]byte, error) {
	sm := re.FindSubmatch(f.text)
	if len(sm) != 6 {
		return nil, errors.New("cannot find elements")
	}

	m := make(map[string][]byte, 5)
	for i, n := range re.SubexpNames() {
		if i > 0 && n != "" {
			m[n] = sm[i]
		}
	}

	return m, nil
}

func (f *Formula) read() (element, error) {
	e := element{}
	text, err := ioutil.ReadFile(f.path)
	if err != nil {
		return e, errors.Wrap(err, "error in ReadFile")
	}

	f.text = text
	m, err := f.findMatches()
	if err != nil {
		return e, errors.Wrap(err, "error in findMatches")
	}

	e.version = m["version"]
	e.tag = m["tag"]
	e.dmg = convBytes(m["dmg"])
	e.appcast = convBytes(m["appcast"])
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

func (f *Formula) save() error {
	repl := []byte(fmt.Sprintf(tmpl, string(f.version), f.tag, f.dmg, f.appcast))
	f.text = re.ReplaceAll(f.text, repl)
	if err := ioutil.WriteFile(f.path, f.text, 0644); err != nil {
		return errors.Wrap(err, "error in WriteFile")
	}
	return nil
}
