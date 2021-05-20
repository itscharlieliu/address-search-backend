package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/itscharlieliu/address-search-backend/utils"
)

func addToConfidenceMap(addresses *csvAddresses, confidenceMap *map[int]int, cleanedString string, addressIdx int, columnIdx int) {
	// If the current column contains a search string, then we increment the confidence of the current address

	if utils.StringContains((*addresses)[addressIdx][columnIdx], (cleanedString)) {
		if confidence, ok := (*confidenceMap)[addressIdx]; ok {
			(*confidenceMap)[addressIdx] = confidence + 1
			return
		}
		(*confidenceMap)[addressIdx] = 1
	}
}

func filterAddresses(addresses *csvAddresses, query string) csvAddresses {
	// Remove special characters from the string and split on spaces
	cleanedStrings := strings.Fields(utils.RemoveSpecialChars(query))

	// Create a map to store how many hits we get for an individual address
	confidenceMap := make(map[int]int)

	// This time complexity is O(N * M)
	// Where N is the number of addresses we have stored, and M is the number of words in the query
	// If performance is a concern, we can limit the number of words allowed
	for i := 0; i < len(*addresses); i++ {
		for j := 0; j < len(cleanedStrings); j++ {
			// The four columns are address, city, state, and zip.
			// This allows us to add more columns in the future if we wish to filter on those
			addToConfidenceMap(addresses, &confidenceMap, cleanedStrings[j], i, 3)
			addToConfidenceMap(addresses, &confidenceMap, cleanedStrings[j], i, 4)
			addToConfidenceMap(addresses, &confidenceMap, cleanedStrings[j], i, 5)
			addToConfidenceMap(addresses, &confidenceMap, cleanedStrings[j], i, 6)
		}
	}

	// Make a slice from the map so that we can sort it by confidence
	confidenceSlice := make([]int, 0, len(confidenceMap))
	for addressIdx := range confidenceMap {
		confidenceSlice = append(confidenceSlice, addressIdx)
	}

	// Sorting the slice by confidence level
	sort.SliceStable(confidenceSlice, func(i, j int) bool {
		return confidenceMap[confidenceSlice[i]] >= confidenceMap[confidenceSlice[j]]
	})

	// Constructing the filtered results ordered by confidence
	filteredResults := make(csvAddresses, 0, len(confidenceSlice))

	for i := range confidenceSlice {
		filteredResults = append(filteredResults, (*addresses)[confidenceSlice[i]])
	}

	return filteredResults
}

// Allows user to call GET on the search endpoint and recieve the resulting items in JSON format
// Customize the search using query parameters
// 	- Note: Only the first instance of the query parameter is used
// Accepted query parameters:
// 	- query: The string that we are using to search all addresses
// We are only processing address params for this exercise, but other params can be easily added
func (handler *BaseHandler) Search(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()

	var results csvAddresses

	// If the query parameter exists, we use that to filter the addresses
	if query, ok := queryParams["query"]; ok {
		results = filterAddresses(handler.addresses, query[0])
	} else {
		results = csvAddresses{} // Empty slice
	}

	// Finally, we respond to the http request
	bytes, err := json.Marshal(csvToStructs(results))

	if err != nil {
		log.Panicln("Unable to generate json: " + err.Error())
	}

	writer.Write([]byte(bytes))
}
