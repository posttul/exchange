package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/posttul/exchange/storage"
)

func main() {
	port := flag.String("port", ":3000", "The port to be monted in.")
	rateDB := flag.String("ratedb", "rates.db", "The default name for the db.")
	authDB := flag.String("authdb", "auth.db", "The default name for the db.")
	flag.Parse()

	// Start storage.
	log.Printf("Set storage to file %s", *rateDB)
	store, err := storage.NewFileStorage(*rateDB)
	if err != nil {
		panic(err)
	}
	// Start Server
	s := &Server{
		Storage: store,
	}

	authstore, err := storage.NewFileStorage(*authDB)
	if err != nil {
		panic(err)
	}
	auth := AuthServer{
		Store: authstore,
	}

	scrap := NewScraper(time.Second * 5)
	go scrap.GetData(store)

	log.Printf("Starting http server a port %s", *port)
	http.HandleFunc("/rate", s.GetUSDRate())
	http.HandleFunc("/gettoken", auth.GetToken())
	http.ListenAndServe(*port, nil)
}
