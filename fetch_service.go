package main

import (
	"encoding/json"
	"errors"
	"log"
)

type fetchArgs struct {
	address string
	id      int
}

type fetchReply struct {
	status bool
	err    string
}

func FetchService(queue string, args ...interface{}) error {
	log.Println(args)

	website, ok := args[0].(string)
	if !ok {
		log.Println("Bad argument website")
		return errors.New("Bad argument website")
	}

	idNum, ok := args[1].(json.Number)
	if !ok {
		log.Println("Bad argument id")
		return errors.New("Bad argument ID")
	}

	id, err := idNum.Int64()

	if err != nil {
		log.Println("bad conversion to int64...")
		return errors.New("bad conversion to int64")
	}

	log.Println("Starting...")

	go Dsptchr.DownloadAndStore(&fetchArgs{
		address: website,
		id:      int(id),
	})

	return nil
}
