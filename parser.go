package main

import (
	"io"

	"golang.org/x/net/html"
	// "strings"
)

type Parser struct{}

func (p *Parser) isRefTag(t *html.Token) bool {
	switch t.Data {
	case "img", "script", "link":
		return true
	}
	return false
}

func (p *Parser) ParseLinks(body io.Reader, website *Website) []*XDep {
	z := html.NewTokenizer(body)

	links := make([]*XDep, 0)

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return links
		case html.SelfClosingTagToken:
			t := z.Token()

			if !p.isRefTag(&t) {
				continue
			}

			xdep, err := NewXDep(&t, website)

			if err == nil {
				links = append(links, xdep)
			}
		}
	}

	return nil
}
