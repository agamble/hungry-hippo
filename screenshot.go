package main

import (
	"bytes"
	"image/png"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type Screenshot struct {
	url         string
	body        []byte
	website     *Website
	thumbnails  []Storable
	isThumbnail bool
	fileName    string
}

const ManetHost string = "http://manet:8891"

func (s *Screenshot) UploadPath() string {
	return filepath.Join(s.website.UploadFolder(), s.Filename())
}

func (s *Screenshot) StoreUrl() string {
	return BASE_STORAGE_URL + filepath.Join(BUCKET_NAME, s.UploadPath())
}

func (s *Screenshot) Filename() string {
	if s.fileName == "" {
		s.fileName = s.genFileName()
	}
	return s.fileName
}

func (s *Screenshot) ScreenshotUrl() string {
	return s.StoreUrl()
}

func (s *Screenshot) ThumbnailUrl() (url string) {
	if s.thumbnails == nil {
		return
	}

	s, ok := s.thumbnails[0].(*Screenshot)
	if !ok {
		return
	}

	url = s.StoreUrl()
	return
}

func (s *Screenshot) genFileName() string {
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

	return gen(16) + s.Ext()
}

func (s *Screenshot) Body() []byte {
	return s.body
}

func (s *Screenshot) PrepareUploadXDeps() {
	if s.isThumbnail {
		return
	}

	img, err := png.Decode(bytes.NewBuffer(s.body))

	if err != nil {
		log.Printf("Error decoding")
		return
	}

	img, err = cutter.Crop(img, cutter.Config{
		Width:  1000,
		Height: 800,
		Mode:   cutter.TopLeft,
	})

	if err != nil {
		log.Printf("Failed cropping image")
		return
	}

	img = resize.Thumbnail(500, 1000, img, resize.Lanczos2)

	if err != nil {
		log.Printf("Failed resizing img")
		return
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)

	if err != nil {
		log.Printf("Error converting to png")
		return
	}

	s.thumbnails = append(s.thumbnails, &Screenshot{
		url:         s.url,
		body:        buf.Bytes(),
		isThumbnail: true,
		website:     s.website,
	})

	return
}

func (s *Screenshot) UploadXDeps() []Storable {
	return s.thumbnails
}

func (s *Screenshot) Ext() string {
	return ".png"
}

func (s *Screenshot) SetStatusCode(code int) {
}

func (s *Screenshot) Url() string {
	u, _ := url.Parse(ManetHost)
	parameters := url.Values{}
	parameters.Add("url", s.url)

	u.RawQuery = parameters.Encode()
	return u.String()
}

func (s *Screenshot) SetBody(b []byte) {
	s.body = b
}

func (s *Screenshot) Dependencies() []Downloadable {
	return nil
}

func (s *Screenshot) PrepareDependencies() {

}

func NewScreenshot(address string, website *Website) *Screenshot {
	s := new(Screenshot)
	s.website = website
	s.url = address
	return s
}
