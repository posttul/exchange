package main

import (
	"encoding/json"
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
		rates, err := s.Storage.Read()
		if err != nil {
			bts, err := json.Marshal(rates)
			if err != nil {
				w.WriteHeader(http.StatusOK)
				w.Write(bts)
			}
		}
	}
}
