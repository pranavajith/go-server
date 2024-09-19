package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type Server struct {
	store PlayerStore
}

// server.go
func (p *Server) processWin(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

// server.go
func (p *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, r)
	case http.MethodGet:
		p.showScore(w, r)
	}

}

func (p *Server) showScore(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

// func Server(w http.ResponseWriter, r *http.Request) {
// 	player := strings.TrimPrefix(r.URL.Path, "/players/")

// 	switch player {
// 	case "Floyd":
// 		fmt.Fprint(w, "10")
// 	case "Pepper":
// 		fmt.Fprint(w, "20")
// 	default:
// 		fmt.Fprint(w, "")
// 	}
// }
