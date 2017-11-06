package exchange

import (
	"coinHelper/protocol"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type bithumb struct {
}

// UseBithumb is a empty instance(?) of bithumb struct for using its methods
var (
	UseBithumb             bithumb
	bithumbCache           = map[string]SearchResult{}
	lastBithumbCacheUpdate time.Time
)

/**
 * Search: search for current currency exchange infomation for given currency
 */
func (bithumb) Search(currencyName string) (SearchResult, error) {
	cachedResult := bithumbCache[currencyName]

	if time.Now().Sub(cachedResult.Time) < minTimeElapse {
		return cachedResult, nil
	}

	currencyCode := CurrencyCodes[currencyName]

	if len(currencyCode) == 0 {
		return SearchResult{}, ErrCoinTypeNotAvailable
	}

	netClient := &http.Client{
		Timeout: timeout,
	}
	response, err := netClient.Get("https://api.bithumb.com/public/ticker/" + currencyCode)

	if err != nil {
		log.Println("Error calling bithumb api: ", err)
		return SearchResult{}, ErrSearchAPI
	}

	var bithumbResp protocol.BithumbTickerResponse
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&bithumbResp); err != nil {
		log.Println("Error decoding bithumb response: ", err)
		return SearchResult{}, ErrJSONParse
	}

	data := bithumbResp.Data

	bithumbCache[currencyName] = SearchResult{
		LastPrice:    int(data.Last),
		Change:       int(data.Last - data.First),
		ChangeRate:   (data.Last - data.First) / data.First * 100,
		CurrencyName: currencyName,
		Time:         time.Now(),
	}

	return bithumbCache[currencyName], nil
}

/**
 * SearchAll: search for current currency exchange infomation for all supported currencies
 */
func (bithumb) SearchAll() ([]SearchResult, error) {
	resultList := []SearchResult{}
	if time.Now().Sub(lastBithumbCacheUpdate) < minTimeElapse {
		for currencyName := range bithumbCache {
			resultList = append(resultList, bithumbCache[currencyName])
		}
		return resultList, nil
	}

	netClient := &http.Client{
		Timeout: timeout,
	}
	response, err := netClient.Get("https://api.bithumb.com/public/ticker/ALL")
	if err != nil {
		log.Println("Error calling bithumb api: ", err)
		return []SearchResult{}, ErrSearchAPI
	}

	var bithumbResp protocol.BithumbTickerAllResponse
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&bithumbResp); err != nil {
		log.Println("Error decoding bithumb response: ", err)
		return []SearchResult{}, ErrJSONParse
	}

	data := bithumbResp.Data

	for currencyName, currencyCode := range CurrencyCodes {
		currencyInfo := data[currencyCode].(map[string]interface{})

		lastPrice, _ := strconv.Atoi(currencyInfo["closing_price"].(string))
		firstPrice, _ := strconv.Atoi(currencyInfo["opening_price"].(string))

		bithumbCache[currencyName] = SearchResult{
			LastPrice:    lastPrice,
			Change:       lastPrice - firstPrice,
			ChangeRate:   float64(lastPrice-firstPrice) / float64(firstPrice) * 100,
			CurrencyName: currencyName,
			Time:         time.Now(),
		}

		lastBithumbCacheUpdate = time.Now()
		resultList = append(resultList, bithumbCache[currencyName])
	}

	return resultList, nil
}
