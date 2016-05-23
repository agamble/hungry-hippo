package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benmanns/goworker"
	_ "github.com/lib/pq"
)

var Dsptchr *Dispatcher = NewDispatcher()

var DbAuth string = "host=db user=%s password=%s dbname=%s sslmode=disable"

func main() {
	DbAuth = fmt.Sprintf(DbAuth,
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE"))

	log.Println(DbAuth)

	time.Sleep(2 * time.Second)
	log.Println("Starting...")

	goworker.Register("Fetch", FetchService)

	go Dsptchr.Dispatch()

	if err := goworker.Work(); err != nil {
		log.Println("Error:", err)
	}
}
