package main

import (
	"errors"
	"math/rand"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type XDep struct {
	Website     *Website
	urlString   string
	token       *html.Token
	body        []byte
	newFileName string
}

var ErrInvalidToken = errors.New("Token passed could not be made into external dependency")

func (xd *XDep) Dependencies() []Downloadable {
	return nil
}

func (xd *XDep) SetStatusCode(code int) {
}

func (xd *XDep) PrepareDependencies() {
}

func (xd *XDep) PrepareUploadXDeps() {
}

func (xd *XDep) UploadXDeps() []Storable {
	return nil
}

func (xd *XDep) Filename() string {
	if xd.newFileName == "" {
		xd.newFileName = xd.genFileName()
	}
	return xd.newFileName
}

func (xd *XDep) UrlString() string {
	return xd.urlString
}

func (xd *XDep) Url() string {
	urlString := xd.urlString

	if strings.HasPrefix(xd.urlString, "//") {
		urlString = "http:" + urlString
	}

	url, _ := url.Parse(urlString)

	return url.String()
}

func (xd *XDep) Id() int {
	return xd.Website.id
}

func (xd *XDep) SetBody(b []byte) {
	xd.body = b
}

func (xd *XDep) Body() []byte {
	return xd.body
}

func (xd *XDep) UploadPath() string {
	return filepath.Join(xd.Website.UploadFolder(), xd.Filename())
}

func (xd *XDep) genFileName() string {
	var src = rand.NewSource(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	gen := func(n int) string {
		b := make([]byte, n)
		// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
		for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}

		return string(b)
	}

	ext := filepath.Ext(xd.Url())

	if ext == "" {
		return gen(16)
	}

	if strings.Contains(ext, "?") {
		ext = ext[:strings.Index(ext, "?")]
	}

	return gen(16) + ext
}

func (xd *XDep) Ext() string {
	return filepath.Ext(xd.newFileName)
}

func (xd *XDep) NewToken() *html.Token {
	t := *xd.token
	for i := range t.Attr {
		attr := &t.Attr[i]
		if attr.Key == "src" {
			attr.Val = xd.Filename()
		}
	}
	return &t
}

func (xd *XDep) ReplaceableStrings() []string {
	tags := make([]string, 0)
	oldType := xd.token.Type

	xd.token.Type = html.SelfClosingTagToken
	tags = append(tags, xd.token.String())

	xd.token.Type = html.StartTagToken
	tags = append(tags, xd.token.String())

	xd.token.Type = oldType
	return tags
}

func (xd *XDep) OldToken() *html.Token {
	return xd.token
}

func NewXDep(token *html.Token, url string, website *Website) *XDep {
	xdep := new(XDep)
	xdep.Website = website
	xdep.token = token
	xdep.urlString = url

	return xdep
}
