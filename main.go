package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/posttul/exchange/storage"
)

func main() {
	port := flag.String("port", ":3000", "The port to be monted in.")
	fileName := flag.String("file", "rates.db", "The default name for the db.")
	flag.Parse()

	// Start storage.
	log.Printf("Set storage to file %s", *fileName)
	store, err := storage.NewFileStorage(*fileName)
	defer store.Close()
	if err != nil {
		panic(err)
	}

	// Start Server
	s := Server{
		Storage: store,
	}
	log.Printf("Starting http server a port %s", *port)
	http.HandleFunc("/rate", s.GetUSDRate())
	http.ListenAndServe(*port, nil)
}
