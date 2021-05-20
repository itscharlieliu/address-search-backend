package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/itscharlieliu/address-search-backend/api"
)

func main() {

	// Read in the csv file
	// TODO add flag parsing to be able to read in different files
	addressesFile, err := os.Open("./redfin_2021-05-19-11-16-39.csv")

	if err != nil {
		log.Fatal("Unable to open file")
	}

	reader := csv.NewReader(addressesFile)

	addresses, err := reader.ReadAll()

	// Cut out the first line. We don't want the headers
	addresses = addresses[1:]

	if err != nil {
		log.Fatal("Unable to parse csv")
	}

	addressesFile.Close()

	// Create the base handler to hold all the data
	handler := api.NewBaseHandler(&addresses)

	// Register the api endpoint for searching
	http.HandleFunc("/search", handler.Search)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
