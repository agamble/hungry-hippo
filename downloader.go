package main

import (
	// "bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type Downloader struct {
	workers  int
	jobsC    chan Downloadable
	ResultsC chan *downloadResult
	Client   *http.Client
}

type Downloadable interface {
	Url() string
	SetBody(b []byte)
	Dependencies() []Downloadable
	PrepareDependencies()
}

type downloadResult struct {
	downloaded Downloadable
	success    bool
}

func (d *Downloader) Download(downloadable Downloadable) {
	resp, err := d.Client.Get(downloadable.Url())
	log.Println("Getting...", downloadable.Url())

	if err != nil {
		log.Println("error getting", err)
		d.ResultsC <- &downloadResult{
			downloaded: nil,
			success:    false,
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Println("error reading")
		d.ResultsC <- &downloadResult{
			downloaded: nil,
			success:    false,
		}
		return
	}

	downloadable.SetBody(body)

	log.Println("Finished fetching...", downloadable.Url())
	d.ResultsC <- &downloadResult{
		downloaded: downloadable,
		success:    true,
	}
}

func NewDownloader() *Downloader {
	d := new(Downloader)
	d.ResultsC = make(chan *downloadResult)
	d.Client = new(http.Client)
	return d
}
