package main

import "net/http"

type APIServer struct {
	listenAddr string
}

func newAPIServer(ListenAddr string) *APIServer {
	return &APIServer{
		listenAddr: ListenAddr,
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
