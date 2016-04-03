package main

import "testing"

func TestDownload(t *testing.T) {
	d := NewDownloader()
	ws, err := NewWebsiteFromAddress("http://www.apple.com", "alex@example.com")
	if err != nil {
		t.Fail()
	}

	go d.Download(ws)
	res := <-d.ResultsC

	if !res.success {
		t.FailNow()
	}

}
