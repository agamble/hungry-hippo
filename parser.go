package main

import (
	"io"
	"log"
	"net/url"

	"golang.org/x/net/html"
	"strings"
)

type Parser struct{}

func (p *Parser) getAttr(t *html.Token, key string) (string, bool) {
	for _, a := range t.Attr {
		if a.Key == key {
			if a.Val == "" {
				return "", false
			}
			return a.Val, true
		}
	}

	return "", false
}

func (p *Parser) fixUrlString(urlString string) string {
	if strings.HasPrefix(urlString, "//") {
		urlString = "http:" + urlString
	}

	return urlString
}

func (p *Parser) isDownloadableToken(t *html.Token) (urlString string, ok bool) {
	switch t.Data {
	case "img", "script":
		urlString, ok = p.getAttr(t, "src")

		if !ok {
			return
		}
	case "link":
		var attr string
		attr, ok = p.getAttr(t, "rel")

		if !ok {
			return
		}

		if attr != "stylesheet" {
			return "", false
		}

		urlString, ok = p.getAttr(t, "href")

		if !ok {
			return
		}
	}

	if urlString == "" {
		return "", false
	}

	return
}

func (p *Parser) ParseLinks(body io.Reader, website *Website) []*XDep {
	z := html.NewTokenizer(body)

	links := make([]*XDep, 0)

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			log.Println("leaving parser")

			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			urlString, ok := p.isDownloadableToken(&t)

			if !ok {
				continue
			}

			fixedUrlString := p.fixUrlString(urlString)
			url, err := url.Parse(fixedUrlString)

			if err != nil {
				continue
			}

			if url.Scheme == "data" {
				continue
			}

			xdep := NewXDep(&t, urlString, website)

			log.Printf("Adding xdep of type %s with url %s", t.Data, urlString)
			links = append(links, xdep)
		}
	}

	return nil
}
