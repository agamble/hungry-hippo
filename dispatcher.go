package main

import "log"

type Dispatcher struct {
	ServiceC   chan *fetchArgs
	Downloader *Downloader
	Uploader   *Uploader

	FinishedC chan interface{}
}

func (d *Dispatcher) DownloadAndStore(args *fetchArgs) {
	d.ServiceC <- args
}

func (d *Dispatcher) DownloadAndStoreWebsite(ws *Website) {
	go d.Downloader.Download(ws)
}

func (d *Dispatcher) Dispatch() {
	downloader := d.Downloader
	uploader := d.Uploader
	uploader.Init()

	for {
		select {
		case args := <-d.ServiceC:
			ws, err := NewWebsiteFromAddress(args.address, args.id)
			if err != nil {
				log.Println("Error building website...")
				continue
			}

			go downloader.Download(ws)
		case downloadResult := <-downloader.ResultsC:
			if !downloadResult.success {
				continue
			}

			downloaded := downloadResult.downloaded
			downloaded.PrepareDependencies()

			for _, xDep := range downloaded.Dependencies() {
				log.Println("dep")
				go downloader.Download(xDep)
			}

			storable, ok := downloaded.(Storable)

			if !ok {
				continue
			}

			go uploader.Upload(storable)

			storable.PrepareUploadXDeps()
			for _, f := range storable.UploadXDeps() {
				go uploader.Upload(f)
			}
		case uploadResult := <-uploader.ResultsC:
			documentable, ok := uploadResult.(Documentable)

			if ok {
				go documentable.SaveReference()
			}

			if d.FinishedC != nil {
				d.FinishedC <- uploadResult
			}
		}

	}
}

func NewDispatcher() *Dispatcher {
	dis := new(Dispatcher)
	dis.Downloader = NewDownloader()
	dis.Uploader = NewUploader()
	dis.ServiceC = make(chan *fetchArgs)

	return dis
}
