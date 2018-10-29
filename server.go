package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/posttul/exchange/storage"
)

// Server is the handler for the http service
type Server struct {
	Storage storage.Storage
}

// GetUSDRate returns the exchange rate for the usd.
func (s *Server) GetUSDRate() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rates, err := s.Storage.Read(&storage.Response{})
		if err != nil {
			log.Printf("The storage is failing, err -> %s", err.Error())
			fail(w, "Service down")
			return
		}
		bts, err := json.Marshal(rates)
		if err != nil {
			log.Printf("Something fail on json GetUSDRate, err -> %s", err.Error())
			fail(w, "internal server error")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bts)
	}
}

func fail(w http.ResponseWriter, say string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"error","msg":"%s"}`, say)))
}
