package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Here are the headers in the csv provided by redfin
// SALE TYPE,SOLD DATE,PROPERTY TYPE,ADDRESS,CITY,STATE OR PROVINCE,ZIP OR POSTAL CODE,PRICE,BEDS,BATHS,LOCATION,SQUARE FEET,LOT SIZE,YEAR BUILT,DAYS ON MARKET,$/SQUARE FEET,HOA/MONTH,STATUS,NEXT OPEN HOUSE START TIME,NEXT OPEN HOUSE END TIME,URL (SEE http://www.redfin.com/buy-a-home/comparative-market-analysis FOR INFO ON PRICING),SOURCE,MLS#,FAVORITE,INTERESTED,LATITUDE,LONGITUDE
type address struct {
	SaleType               string
	SoldDate               string
	PropertyType           string
	Address                string
	City                   string
	StateOrProvince        string
	ZipOrPostalCode        string
	Price                  string
	Beds                   string
	Baths                  string
	Location               string
	SquareFeet             string
	LotSize                string
	YearBuilt              string
	DaysOnMarket           string
	PricePerSquareFeet     string
	HoaPerMonth            string
	NextOpenHouseStartTime string
	NextOpenHouseEndTime   string
	Url                    string
	Source                 string
	MlsNumber              string
	Favorite               string
	Interested             string
	Latitude               string
	Longitude              string
}

func csvToStructs(addresses [][]string) (results []address) {

	for i := 0; i < len(addresses); i++ {
		results = append(results, address{
			SaleType:               addresses[i][0],
			SoldDate:               addresses[i][1],
			PropertyType:           addresses[i][2],
			Address:                addresses[i][3],
			City:                   addresses[i][4],
			StateOrProvince:        addresses[i][5],
			ZipOrPostalCode:        addresses[i][6],
			Price:                  addresses[i][7],
			Beds:                   addresses[i][8],
			Baths:                  addresses[i][9],
			Location:               addresses[i][10],
			SquareFeet:             addresses[i][11],
			LotSize:                addresses[i][12],
			YearBuilt:              addresses[i][13],
			DaysOnMarket:           addresses[i][14],
			PricePerSquareFeet:     addresses[i][15],
			HoaPerMonth:            addresses[i][16],
			NextOpenHouseStartTime: addresses[i][17],
			NextOpenHouseEndTime:   addresses[i][18],
			Url:                    addresses[i][19],
			Source:                 addresses[i][20],
			MlsNumber:              addresses[i][21],
			Favorite:               addresses[i][22],
			Interested:             addresses[i][23],
			Latitude:               addresses[i][24],
			Longitude:              addresses[i][25],
		})
	}

	return results
}

// Allows user to call GET on the search endpoint and recieve the resulting items in JSON format
// Customize the search using query parameters
// Note - Only the first instance of the query parameter is used
// Accepted query parameters:
// 	- address: Return any items that include the specified address
func (handler *BaseHandler) Search(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()

	// Make a copy of the addresses for us to work with without worrying about mutating the original
	result := make([][]string, len(*handler.addresses))
	copy(result, *handler.addresses)

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
