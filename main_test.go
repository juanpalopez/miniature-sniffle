package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() *httptest.Server {
	mu.Lock()
	animes = make(map[int]*Anime)
	byTitle = make(map[string]int)
	nextID = 1
	mu.Unlock()
	return httptest.NewServer(newServer())
}

func TestCreateAnime(t *testing.T) {
	ts := setup()
	defer ts.Close()

	body := `{"title":"Naruto","creator":"Masashi","genre":"Shonen","year":2002}`
	resp, err := http.Post(ts.URL+"/animes", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", resp.StatusCode)
	}

	// duplicate title
	resp2, err := http.Post(ts.URL+"/animes", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("post2: %v", err)
	}
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 got %d", resp2.StatusCode)
	}
}

func TestAddReview(t *testing.T) {
	ts := setup()
	defer ts.Close()

	body := `{"title":"Naruto","creator":"Masashi","genre":"Shonen","year":2002}`
	resp, err := http.Post(ts.URL+"/animes", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", resp.StatusCode)
	}

	review := `{"rating":5,"text":"Great"}`
	resp, err = http.Post(ts.URL+"/animes/1/reviews", "application/json", bytes.NewBufferString(review))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", resp.StatusCode)
	}

	r, err := http.Get(ts.URL + "/animes/1")
	if err != nil {
		t.Fatal(err)
	}
	var a Anime
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		t.Fatal(err)
	}
	if len(a.Reviews) != 1 {
		t.Fatalf("expected 1 review got %d", len(a.Reviews))
	}
}
