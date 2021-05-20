package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/itscharlieliu/address-search-backend/utils"
)

func filterResults(results [][]string, queryParams url.Values, queryParam string, paramIdx int) [][]string {
	if queriesAddresses, ok := (queryParams[queryParam]); ok {
		// We make the filteredResult list max size to be the same as the current results size
		filteredResults := [][]string{}
		for i := 0; i < len(results); i++ {
			if utils.StringContains(results[i][paramIdx], (queriesAddresses[0])) {

				filteredResults = append(filteredResults, results[i])
			}
		}
		return filteredResults
	}
	return results
}

// Allows user to call GET on the search endpoint and recieve the resulting items in JSON format
// Customize the search using query parameters
// 	- Note: Only the first instance of the query parameter is used
// Accepted query parameters:
// 	- address: Street address
//	- city: Query city ("San Francisco")
//	- state: 2 character state code ("CA")
//	- zip: ZIP code
// We are only processing address params for this exercise, but other params can be easily added
func (handler *BaseHandler) Search(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()

	// Make a copy of the addresses for us to work with without worrying about mutating the original
	results := make([][]string, len(*handler.addresses))
	copy(results, *handler.addresses)

	// For each query param, we whittle down the size of the results.
	// This approach gives us a time complexity of O(N * M), where M is the number of query parameters implemented.
	// However, because each pass of the filter will whittle down the results, the amortized runtime is much better.

	results = filterResults(results, queryParams, "address", 3)
	results = filterResults(results, queryParams, "city", 4)
	results = filterResults(results, queryParams, "state", 5)
	results = filterResults(results, queryParams, "zip", 6)

	bytes, err := json.Marshal(csvToStructs(results))

	if err != nil {
		log.Panicln("Unable to generate json: " + err.Error())
	}

	writer.Write([]byte(bytes))
}
