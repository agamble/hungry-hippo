package main

import (
	"testing"
	"time"
)

func TestFullDispatch(t *testing.T) {
	ws, _ := NewWebsiteFromAddress("http://bbc.com", "alex@example.com")
	d := NewDispatcher()
	InitDsClient()

	d.FinishedC = make(chan interface{})
	go d.Dispatch()
	d.DownloadAndStoreWebsite(ws)

	timeout := time.After(20 * time.Second)

	select {
	case <-timeout:
		t.FailNow()
	case <-d.FinishedC:
		return
	}

}
