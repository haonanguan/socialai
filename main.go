package main

import (
	"fmt"
	"log"
	"net/http"
	"socialai/backend"
	"socialai/handler"
)

func main() {
	fmt.Println("started-service")

	backend.InitElasticsearchBackend()
	backend.InitGCSBackend()

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
