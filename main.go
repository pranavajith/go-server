package main

// func main() {
// 	server := NewServer(":4040")
// 	server.Run()
// }

import (
	"log"
	"net/http"
)

func main() {
	server := &Server{&InMemoryPlayerStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
