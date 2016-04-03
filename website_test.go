package main

import "testing"

func MockWebsite() *Website {
	ws, err := NewWebsiteFromAddress("http://apple.com", "alex@example.com")

	if err != nil {
		panic(err)
	}

	return ws
}

func TestDataStoreUrl(t *testing.T) {
	ws := MockWebsite()
	InitDsClient()
	err := ws.SaveReference()

	if err != nil {
		panic(err)
		t.FailNow()
	}
}
