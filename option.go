package gumk

import (
	"context"
	"fmt"
	"net/http"
)

// Option means an interface of options
type Option interface {
	apply(*Gumk)
}

// WithContext is an option to determine context
func WithContext(ctx context.Context) Option { return withContext{ctx} }

type withContext struct{ ctx context.Context }

func (w withContext) apply(g *Gumk) { g.context = w.ctx }

// WithTag is an option to determine the tag to download
func WithTag(tag string) Option { return withTag{tag} }

type withTag struct{ tag string }

func (w withTag) apply(g *Gumk) { g.tag = w.tag }

// WithFormula is an option to determine the path of the formula
func WithFormula() Option {
	return withFormula{"/usr/local/Homebrew/Library/Taps/delphinus/homebrew-macvim-kaoriya/Casks/macvim-kaoriya.rb"}
}

type withFormula struct{ path string }

func (w withFormula) apply(g *Gumk) { g.formula = w.path }

// WithDMG is an option to determine the URL of the DMG file
func WithDMG(tag string) Option {
	return withDMG{
		"https://github.com/splhack/macvim-kaoriya/releases/download/%s/MacVim-KaoriYa-%s.dmg",
		tag,
	}
}

type withDMG struct{ url, tag string }

func (w withDMG) apply(g *Gumk) { g.dmg = func() string { return fmt.Sprintf(w.url, w.tag, w.tag) } }

// WithRelease is an option to determine the URL of the release file
func WithRelease(tag string) Option {
	return withRelease{
		"https://github.com/splhack/macvim-kaoriya/releases/tag/%s",
		tag,
	}
}

type withRelease struct{ url, tag string }

func (w withRelease) apply(g *Gumk) { g.release = func() string { return fmt.Sprintf(w.url, w.tag) } }

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
