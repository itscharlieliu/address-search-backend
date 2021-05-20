package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Allows user to call GET on the search endpoint and recieve the resulting items in JSON format
// Customize the search using query parameters
// 	- Note: Only the first instance of the query parameter is used
// Accepted query parameters:
// 	- address: Return any items that include the specified address
func (handler *BaseHandler) Search(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()

	// Make a copy of the addresses for us to work with without worrying about mutating the original
	result := make([][]string, len(*handler.addresses))
	copy(result, *handler.addresses)

	// For each query param, we whittle down the size of the results.
	// This approach gives us a time complexity of O(N * M), where M is the number of query parameters implemented.
	// However, because each pass of the filter will whittle down the results, the amortized runtime is much better.
	if queriesAddresses, ok := (queryParams["address"]); ok {
		// We make the filteredResult list max size to be the same as the current results size
		filteredResult := [][]string{}
		for i := 0; i < len(result); i++ {
			if strings.Contains(result[i][3], (queriesAddresses[0])) {

				filteredResult = append(filteredResult, result[i])
			}
		}
		result = filteredResult
	}

	bytes, err := json.Marshal(csvToStructs(result))

	if err != nil {
		log.Panicln("Unable to generate json: " + err.Error())
	}

	writer.Write([]byte(bytes))
}
