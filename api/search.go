package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	lru "github.com/hashicorp/golang-lru"
	"github.com/itscharlieliu/address-search-backend/utils"
)

func addToConfidenceMap(confidenceMap *map[int]int, addressIdx int) {
	if confidence, ok := (*confidenceMap)[addressIdx]; ok {
		(*confidenceMap)[addressIdx] = confidence + 1
		return
	}
	(*confidenceMap)[addressIdx] = 1
}

// Check if the row contains the requested word. If it does, then increment the confidence of that row in the confidence map
// Returns if the row contains the word
func checkAddressForString(addresses *csvAddresses, confidenceMap *map[int]int, cleanedString string, addressIdx int, columnIdx int) bool {
	// If the current column contains a search string, then we increment the confidence of the current address

	if utils.StringContains((*addresses)[addressIdx][columnIdx], (cleanedString)) {

		addToConfidenceMap(confidenceMap, addressIdx)
		return true
	}
	return false
}

func filterAddresses(addresses *csvAddresses, searchCache *lru.Cache, query string) csvAddresses {
	// Remove special characters from the string and split on spaces
	cleanedStrings := strings.Fields(utils.RemoveSpecialChars(query))

	// Create a map to store how many hits we get for an individual address
	confidenceMap := make(map[int]int)

	// This time complexity is O(N * M)
	// Where N is the number of addresses we have stored, and M is the number of words in the query
	// If performance is a concern, we can limit the number of words allowed

	for j := 0; j < len(cleanedStrings); j++ {
		// If our word is in the cache, then we just need to look through that and add it to the confidence map
		// instead of going through the entire collection
		if addressesFromCache, ok := searchCache.Get(cleanedStrings[j]); ok {
			fmt.Println("Contains" + cleanedStrings[j])

			fmt.Println(addressesFromCache)

			typedAddressesFromCache := addressesFromCache.([]int)

			for i := range typedAddressesFromCache {
				addToConfidenceMap(&confidenceMap, typedAddressesFromCache[i])
			}
			continue
		}

		occurances := []int{}

		for i := 0; i < len(*addresses); i++ {
			// The four columns are address, city, state, and zip.
			// This allows us to add more columns in the future if we wish to filter on those
			if checkAddressForString(addresses, &confidenceMap, cleanedStrings[j], i, 3) ||
				checkAddressForString(addresses, &confidenceMap, cleanedStrings[j], i, 4) ||
				checkAddressForString(addresses, &confidenceMap, cleanedStrings[j], i, 5) ||
				checkAddressForString(addresses, &confidenceMap, cleanedStrings[j], i, 6) {
				// If any of these are true, we know that the word was found in this row.
				occurances = append(occurances, i)
			}
		}
		searchCache.Add(cleanedStrings[j], occurances)
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
		log.Printf("Recieved query: %s\n", query)
		results = filterAddresses(handler.addresses, handler.searchCache, query[0])
	} else {
		results = csvAddresses{} // Empty slice
	}

	// Finally, we respond to the http request
	bytes, err := json.Marshal(csvToStructs(results))

	if err != nil {
		log.Panicln("Unable to generate json: " + err.Error())
	}

	writer.Header().Set("Content-Type", "application/json")
	// Currently allow all origins, but should be limited in a production server
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Write([]byte(bytes))
}
