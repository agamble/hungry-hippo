package main

import (
	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
)

var DatastoreClient *datastore.Client

type Documentable interface {
	SaveReference() error
}

func InitDsClient() {
	ds, err := datastore.NewClient(context.TODO(), "treasure-1270")
	if err != nil {
		panic(err)
	}
	DatastoreClient = ds
}
