package main

// var DatastoreClient *datastore.Client

type Documentable interface {
	SaveReference() error
}
