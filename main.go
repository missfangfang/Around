package main

import (
	"fmt"
	"log"
	"net/http"

	"around/backend"
	"around/handler"
)

func main() {
	fmt.Println("started-service") // helps to debug
	backend.InitElasticsearchBackend()
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
