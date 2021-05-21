package main

import (
	"encoding/csv"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/itscharlieliu/address-search-backend/api"
)

const DEFAULT_FILENAME = "./redfin_2021-05-19-11-16-39.csv"

func main() {

	filenamePtr := flag.String("f", DEFAULT_FILENAME, "Filename of the address csv")
	helpPtr := flag.Bool("h", false, "Display help information")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		return
	}

	// Read in the csv file
	addressesFile, err := os.Open(*filenamePtr)

	if err != nil {
		log.Fatalf("Unable to open file: %s", *filenamePtr)
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
	// Currently will serve at http://localhost:8080/search
	http.HandleFunc("/search", handler.Search)

	log.Println("Listening...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
