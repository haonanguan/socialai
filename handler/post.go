package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"socialai/backend"
	"socialai/model"
	"socialai/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. process http request: multipart -> post struct + file
	fmt.Println("Received one upload request")

	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	p := model.Post{
		Id:      uuid.New(),
		User:    username.(string),
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename)
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	//2. call service level to handle business logic
	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	//3. response
	fmt.Println("Post is saved successfully.")
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
		fmt.Println("Search by user: " + user)
		posts, err = service.SearchPostsByUser(user)
	} else {
		fmt.Println("Search by keywords: " + keywords)
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for delete")

	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	id := mux.Vars(r)["id"]

	if err := service.DeletePost(id, username); err != nil {
		http.Error(w, "Failed to delete post: ", http.StatusInternalServerError)
		fmt.Printf("Failed to delete post: %v\n", err)
		return
	}

	if err := backend.GCSBackend.DeleteFromGCS(id); err != nil {
		http.Error(w, "Failed to delete post from GCS", http.StatusInternalServerError)
		fmt.Printf("Failed to delete post from GCS %v\n", err)
		return
	}

	fmt.Println("Post is deleted successfully")
}
