package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Review struct {
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

type Anime struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Creator string   `json:"creator"`
	Genre   string   `json:"genre"`
	Year    int      `json:"year"`
	Reviews []Review `json:"reviews,omitempty"`
}

var (
	mu      sync.Mutex
	animes  = make(map[int]*Anime)
	byTitle = make(map[string]int)
	nextID  = 1
)

func newServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("%s /animes", http.MethodPost), createAnime)
	mux.HandleFunc(fmt.Sprintf("%s /animes", http.MethodGet), listAnimes)
	mux.HandleFunc(fmt.Sprintf("%s /animes/{id}", http.MethodGet), getAnime)
	mux.HandleFunc(fmt.Sprintf("%s /animes/{id}/reviews", http.MethodPost), addReview)
	mux.HandleFunc(fmt.Sprintf("%s /animes/{id}/reviews", http.MethodGet), listReviews)
	return mux
}

func main() {
	mux := newServer()

	// TODO: add authentication and validation middleware

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func createAnime(w http.ResponseWriter, r *http.Request) {
	var a Anime
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if a.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, ok := byTitle[a.Title]; ok {
		http.Error(w, "anime title already exists", http.StatusConflict)
		return
	}
	a.ID = nextID
	nextID++
	animes[a.ID] = &a
	byTitle[a.Title] = a.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func listAnimes(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	list := make([]*Anime, 0, len(animes))
	for _, a := range animes {
		list = append(list, a)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func getAnime(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	mu.Lock()
	a, ok := animes[id]
	mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func addReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	mu.Lock()
	a, ok := animes[id]
	mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}

	var rv Review
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if rv.Rating < 1 || rv.Rating > 5 {
		http.Error(w, "rating must be between 1 and 5", http.StatusBadRequest)
		return
	}
	mu.Lock()
	a.Reviews = append(a.Reviews, rv)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func listReviews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	mu.Lock()
	a, ok := animes[id]
	if !ok {
		mu.Unlock()
		http.NotFound(w, r)
		return
	}
	reviews := append([]Review(nil), a.Reviews...)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
