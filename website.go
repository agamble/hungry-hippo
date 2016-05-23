package main

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"time"
)

type Website struct {
	url        *url.URL
	dirName    string
	id         int
	statusCode int
	storeUrl   string
	screenshot *Screenshot

	body []byte

	xDeps []Downloadable
}

type Revision struct {
	Status    string
	Owner     int
	Url       string
	StoreUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	FolderPrefix = "wb"

	HtmlFilename = "index.html"

	Downloaded  = "DL"
	Initialised = "IN"
)

func (w *Website) ResolveReference(ref *url.URL) *url.URL {
	return w.url.ResolveReference(ref)
}

func (w *Website) Url() string {
	return w.url.String()
}

func (w *Website) StoreUrl() string {
	return BASE_STORAGE_URL + filepath.Join(BUCKET_NAME, FolderPrefix, w.dirName, HtmlFilename)
}

func (w *Website) SetBody(b []byte) {
	w.body = b
}

func (w *Website) PrepareUploadXDeps() {}

func (w *Website) UploadXDeps() []Storable {
	return nil
}

func (w *Website) Dependencies() []Downloadable {
	return w.xDeps
}

func (w *Website) Body() []byte {
	return w.body
}

func (w *Website) Id() int {
	return w.id
}

func (w *Website) UploadFolder() string {
	return filepath.Join(FolderPrefix, w.dirName)
}

func (w *Website) UploadPath() string {
	return filepath.Join(w.UploadFolder(), HtmlFilename)
}

func (w *Website) Ext() string {
	return ".html"
}

func (w *Website) PrepareDependencies() {
	p := Parser{}
	xDeps := p.ParseLinks(bytes.NewReader(w.body), w)

	log.Println(len(xDeps))

	w.xDeps = make([]Downloadable, len(xDeps)+1)

	for i, xd := range xDeps {
		oldUrl := xd.UrlString()
		newUrl := xd.Filename()

		w.body = bytes.Replace(w.body, []byte(oldUrl), []byte(newUrl), -1)

		w.xDeps[i] = xDeps[i]
	}

	w.xDeps[len(xDeps)] = w.screenshot
}

func (w *Website) SetStatusCode(code int) {
	w.statusCode = code
}

func (w *Website) SaveReference() error {
	db, err := sql.Open("postgres", DbAuth)
	defer db.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = db.Exec("INSERT INTO revisions (page_id, status_code, store_url, thumbnail_url, screenshot_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		w.id,
		w.statusCode,
		w.StoreUrl(),
		w.screenshot.ThumbnailUrl(),
		w.screenshot.ScreenshotUrl(),
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (w *Website) genRandomDirName() {
	randBytes := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	rand.Read(randBytes)
	w.dirName = hex.EncodeToString(randBytes)
}

func NewWebsiteFromAddress(address string, id int) (*Website, error) {
	w := new(Website)
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	w.url = u
	w.id = id

	// make a forbidden page if not allowed
	if err != nil {
		return nil, err
	}

	w.screenshot = NewScreenshot(address, w)

	w.genRandomDirName()
	return w, nil
}
