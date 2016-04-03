package main

import (
	"net/http"
	"net/url"
)

type FetchService struct {
}

type fetchArgs struct {
	Address   string
	UserEmail string
}

type fetchReply struct {
	status bool
	err    string
}

func (f *fetchReply) Fail(message string) {
	f.status = false
	f.err = message
}

func (f *fetchReply) Success() {
	f.status = true
}

func (f *FetchService) Fetch(r *http.Request, args *fetchArgs, reply *fetchReply) {
	_, err := url.Parse(args.Address)

	if err != nil {
		reply.Fail("URL is invalid")
		return
	}

	if args.Address == "" || args.UserEmail == "" {
		reply.Fail("You must include both email and web address as part of request...")
		return
	}

	go Dsptchr.DownloadAndStore(args)
	reply.Success()
}
