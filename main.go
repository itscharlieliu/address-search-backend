package main

import (
	"log"
	"net/http"

	"github.com/itscharlieliu/address-search-backend/api"
)

func main() {
	http.HandleFunc("/", api.Search)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
