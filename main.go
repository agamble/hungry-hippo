package main

import (
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

var Dsptchr *Dispatcher = NewDispatcher()

func main() {
	InitDsClient()

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(FetchService), "")
	http.Handle("/rpc", s)
	http.ListenAndServe(":8080", s)
}
