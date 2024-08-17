package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"socialai/model"
	"socialai/service"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. process http request: json string -> post struct
	fmt.Println("Received one upload request")
	decoder := json.NewDecoder(r.Body)
	var p model.Post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}
	//2. call service level to handle business logic

	//3. response
	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	//1. process http request: URL -> string
	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	//2. call service level to handle business logic
	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	//3. response
	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}
