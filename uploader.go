package main

import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
	"log"
)

type Uploader struct {
	ResultsC chan Storable
}

const (
	BUCKET_NAME      = "hippo"
	BASE_STORAGE_URL = "https://storage.googleapis.com/"
)

type Storable interface {
	UploadPath() string
	Owner() string
	Body() []byte
}

var StorageClient *storage.Client

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
