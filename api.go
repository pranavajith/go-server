package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type Server struct {
	serverAdd string
	posts     (map[int]Post)
	nextId    int
	postsMu   sync.Mutex
}

func NewServer(ServerAdd string) *Server {
	return &Server{
		serverAdd: ServerAdd,
		posts:     make(map[int]Post), // Initialize the map here
		nextId:    1,                  // Start IDs from 1
	}
}

func (s *Server) Run() {
	http.HandleFunc("/posts", s.postsHandler)
	http.HandleFunc("/posts/", s.postHandler)

	fmt.Println("Server running at localhost", s.serverAdd)
	log.Fatal(http.ListenAndServe(s.serverAdd, nil))
}

func (s *Server) postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGetPosts(w)
	case "POST":
		s.handlePostPosts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/posts/"):])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		s.handleGetPost(w, id)
	case "DELETE":
		s.handleDeletePost(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetPosts(w http.ResponseWriter) {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	ps := make([]Post, 0, len(s.posts))
	for _, p := range s.posts {
		ps = append(ps, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ps)
}

func (s *Server) handlePostPosts(w http.ResponseWriter, r *http.Request) {
	var p Post

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	p.ID = s.nextId
	s.nextId++
	s.posts[p.ID] = p

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (s *Server) handleGetPost(w http.ResponseWriter, id int) {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	p, ok := s.posts[id]
	if !ok {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *Server) handleDeletePost(w http.ResponseWriter, id int) {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()
	_, ok := s.posts[id]
	if !ok {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	delete(s.posts, id)
	w.WriteHeader(http.StatusOK)
}
