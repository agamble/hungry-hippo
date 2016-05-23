package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
)

type Uploader struct {
	ResultsC chan Storable
}

const (
	BUCKET_NAME      = "lockbox-elephant"
	BASE_STORAGE_URL = "https://storage.googleapis.com/"
)

type Storable interface {
	UploadPath() string
	Body() []byte
	Ext() string
	SetStatusCode(code int)
	PrepareUploadXDeps()
	UploadXDeps() []Storable
}

var StorageClient *storage.Client

func (u *Uploader) setContentType(file Storable, obj *storage.ObjectHandle) {
	var contentType string

	switch file.Ext() {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "text/javascript"
	default:
		contentType = ""
	}

	attrs := storage.ObjectAttrs{}
	attrs.ContentType = contentType
	obj.Update(context.TODO(), attrs)
	log.Printf("Setting %s, %s with content-type %s", file.UploadPath(), file.Ext(), contentType)
}

func (u *Uploader) Init() {
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	StorageClient = client
}

func (u *Uploader) Upload(file Storable) {
	bucket := StorageClient.Bucket(BUCKET_NAME)
	obj := bucket.Object(file.UploadPath())

	u.writeFile(obj, file)
	u.setAsPublic(obj)
	u.setContentType(file, obj)

	log.Println("Finished storing...", file.UploadPath())
	u.ResultsC <- file
	log.Println("Sent on channel", file.UploadPath())
}

func (u *Uploader) writeFile(obj *storage.ObjectHandle, file Storable) {
	writer := obj.NewWriter(context.TODO())
	log.Println("Writing file...", file.UploadPath())
	writer.Write(file.Body())
	writer.Close()
}

func (u *Uploader) setAsPublic(obj *storage.ObjectHandle) {
	obj.ACL().Set(context.TODO(), storage.AllUsers, storage.RoleReader)
}

func NewUploader() *Uploader {
	u := new(Uploader)
	u.ResultsC = make(chan Storable)
	return u
}
