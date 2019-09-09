package ads

import (
	"fmt"
	"net/http"
	"os"

	"github.com/brave-intl/bat-go/utils/httprequest"
)

var (
	adsAvailable = []string{}
)

type country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// AvailableIn
func AvailableIn() ([]string, error) {
	if len(adsAvailable) != 0 {
		return adsAvailable, nil
	}
	origin := os.Getenv("ADS_URL")
	resp, err := http.Get(origin + "/v1/geoCode")
	if err != nil {
		return adsAvailable, err
	}
	var countries []country
	err = httprequest.ReadBody(resp.Body, &countries)
	if err != nil {
		return adsAvailable, err
	}
	countryList := []string{}
	fmt.Println(countries)
	for _, item := range countries {
		countryList = append(countryList, item.Code)
	}
	// cache the result so we don't ddos ourselves
	adsAvailable = countryList
	return adsAvailable, nil
}
