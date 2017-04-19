package gumk

import (
	"fmt"
	"net/http"
)

// Option means an interface of options
type Option interface {
	apply(*Gumk)
}

// WithFormula is an option to determine the path of the formula
func WithFormula() Option {
	return withFormula{"/usr/local/Homebrew/Library/Taps/delphinus/homebrew-macvim-kaoriya/Casks/macvim-kaoriya.rb"}
}

type withFormula struct{ path string }

func (w withFormula) apply(g *Gumk) { g.formula = w.path }

// WithDMG is an option to determine the URL of the DMG file
func WithDMG(tag string) Option {
	return withDMG{
		fmt.Sprintf("https://github.com/splhack/macvim-kaoriya/releases/download/%s/MacVim-KaoriYa-%s.dmg", tag, tag),
	}
}

type withDMG struct{ url string }

func (w withDMG) apply(g *Gumk) { g.dmg = w.url }

// WithRelease is an option to determine the URL of the release file
func WithRelease(tag string) Option {
	return withRelease{
		fmt.Sprintf("https://github.com/splhack/macvim-kaoriya/releases/tag/%s", tag),
	}
}

type withRelease struct{ url string }

func (w withRelease) apply(g *Gumk) { g.release = w.url }

// WithAppcast is an option to determine the URL of the appcast file
func WithAppcast() Option {
	return withAppcast{"https://github.com/splhack/macvim-kaoriya/releases.atom"}
}

type withAppcast struct{ url string }

func (w withAppcast) apply(g *Gumk) { g.appcast = w.url }

// WithHTTPClient is an option to set HTTP client for the app
func WithHTTPClient(client *http.Client) Option { return withHTTPClient{client} }

type withHTTPClient struct{ client *http.Client }

func (w withHTTPClient) apply(g *Gumk) { g.httpClient = w.client }
