package api

import (
	"encoding/csv"
	"log"
	"os"
	"testing"

	"github.com/itscharlieliu/address-search-backend/utils"
)

func createHandler() BaseHandler {
	addressesFile, err := os.Open("../redfin_2021-05-19-11-16-39.csv")

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
	return *NewBaseHandler(&addresses)
}

func TestFilterAddresses(t *testing.T) {
	handler := createHandler()

	results := filterAddresses(handler.addresses, "")

	if len(results) != 0 {
		t.Error("Filter should not return any results with an empty string")
	}

	results = filterAddresses(handler.addresses, "invalid address")

	if len(results) != 0 {
		t.Errorf("Filtered addressed expected 0; got %d", len(results))
	}

	results = filterAddresses(handler.addresses, "Ave")

	numAddresses := 0
	for i := 0; i < len(*handler.addresses); i++ {
		if utils.StringContains((*handler.addresses)[i][3], "Ave") {
			numAddresses++
		}
	}

	if len(results) != numAddresses {
		t.Errorf("Number of filtered addresses should be %d; got %d", numAddresses, len(results))
	}

}
