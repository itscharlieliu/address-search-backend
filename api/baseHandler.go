package api

// Here are the headers in the csv provided by redfin
// SALE TYPE,SOLD DATE,PROPERTY TYPE,ADDRESS,CITY,STATE OR PROVINCE,ZIP OR POSTAL CODE,PRICE,BEDS,BATHS,LOCATION,SQUARE FEET,LOT SIZE,YEAR BUILT,DAYS ON MARKET,$/SQUARE FEET,HOA/MONTH,STATUS,NEXT OPEN HOUSE START TIME,NEXT OPEN HOUSE END TIME,URL (SEE http://www.redfin.com/buy-a-home/comparative-market-analysis FOR INFO ON PRICING),SOURCE,MLS#,FAVORITE,INTERESTED,LATITUDE,LONGITUDE
// 	- Note: This is assuming the csv format provided by redfin does not change
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
	Status                 string
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

type csvAddresses = [][]string

func csvToStructs(addresses [][]string) []address {
	results := []address{}
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
			Status:                 addresses[i][17],
			NextOpenHouseStartTime: addresses[i][18],
			NextOpenHouseEndTime:   addresses[i][19],
			Url:                    addresses[i][20],
			Source:                 addresses[i][21],
			MlsNumber:              addresses[i][22],
			Favorite:               addresses[i][23],
			Interested:             addresses[i][24],
			Latitude:               addresses[i][25],
			Longitude:              addresses[i][26],
		})
	}

	return results
}

type BaseHandler struct {
	addresses *[][]string // Arrays are already pointers. We don't need to pass in pointer here
}

func NewBaseHandler(addresses *[][]string) *BaseHandler {
	return &BaseHandler{
		addresses: addresses,
	}
}
