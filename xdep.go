package main

import (
	"errors"
	"math/rand"
	"net/url"
	"path/filepath"
	"time"

	"golang.org/x/net/html"
)

type XDep struct {
	Website     *Website
	token       *html.Token
	body        []byte
	newFileName string
}

var ErrInvalidToken = errors.New("Token passed could not be made into external dependency")

func (xd *XDep) Dependencies() []Downloadable {
	return nil
}

func (xd *XDep) PrepareDependencies() {
}

func (xd *XDep) Filename() string {
	if xd.newFileName == "" {
		xd.newFileName = xd.genFileName()
	}
	return xd.newFileName
}

func (xd *XDep) Url() string {
	src, _ := xd.getSrc()
	return src
}

func (xd *XDep) Owner() string {
	return xd.Website.ownerEmail
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

	return gen(16) + ext
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

func (xd *XDep) OldToken() *html.Token {
	return xd.token
}

func (xd *XDep) getSrc() (string, bool) {
	t := xd.token
	for _, a := range t.Attr {
		if a.Key == "src" {
			return a.Val, true
		}
	}

	return "", false
}

func (xd *XDep) tokenOk() bool {
	u, ok := xd.getSrc()

	if !ok {
		return false
	}

	_, err := url.Parse(u)
	if err != nil {
		return false
	}

	return true
}

func NewXDep(token *html.Token, website *Website) (*XDep, error) {
	xdep := new(XDep)
	xdep.Website = website
	xdep.token = token

	if !xdep.tokenOk() {
		return nil, ErrInvalidToken
	}

	return xdep, nil
}
