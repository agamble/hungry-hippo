package main

import (
	"bytes"
	"encoding/hex"
	"log"
	"math/rand"
	"net/url"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
)

type Website struct {
	url     *url.URL
	dirName string

	ownerEmail string

	body []byte

	xDeps []Downloadable
}

type WebsiteStore struct {
	Status    string
	Owner     string
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

func (w *Website) Dependencies() []Downloadable {
	return w.xDeps
}

func (w *Website) Body() []byte {
	return w.body
}

func (w *Website) Owner() string {
	return w.ownerEmail
}

func (w *Website) UploadFolder() string {
	return filepath.Join(FolderPrefix, w.dirName)
}

func (w *Website) UploadPath() string {
	return filepath.Join(w.UploadFolder(), HtmlFilename)
}

func (w *Website) PrepareDependencies() {
	p := Parser{}
	xDeps := p.ParseLinks(bytes.NewReader(w.body), w)

	w.xDeps = make([]Downloadable, len(xDeps))

	for i, xd := range xDeps {
		oldTag := xd.OldToken().String()
		newTag := xd.NewToken().String()

		w.xDeps[i] = xDeps[i]

		w.body = bytes.Replace(w.body, []byte(oldTag), []byte(newTag), -1)
		xd.Website = w
	}
}

func (w *Website) SaveReference() error {
	var status string
	if w.body == nil {
		status = Downloaded
	} else {
		status = Initialised
	}

	wss := &WebsiteStore{
		Status:    status,
		Owner:     w.Owner(),
		Url:       w.Url(),
		StoreUrl:  w.StoreUrl(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	key := datastore.NewIncompleteKey(context.Background(), "website", nil)
	key, err := DatastoreClient.Put(context.Background(), key, wss)

	log.Println(key)

	if err != nil {
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

func NewWebsiteFromAddress(address string, email string) (*Website, error) {
	w := new(Website)
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	w.url = u
	w.ownerEmail = email

	// make a forbidden page if not allowed
	if err != nil {
		return nil, err
	}

	w.genRandomDirName()
	return w, nil
}
