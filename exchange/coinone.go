package exchange

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"coinHelper/protocol"
)

type coinone struct {
}

// UseCoinone is a empty instance(?) of coinone struct for using its methods
var (
	UseCoinone             coinone
	coinoneCache           = map[string]SearchResult{}
	lastCoinoneCacheUpdate time.Time
)

/**
 * Search: search for current currency exchange infomation for given currency
 */
func (coinone) Search(currencyName string) (SearchResult, error) {
	cachedResult := coinoneCache[currencyName]

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
	response, err := netClient.Get("https://api.coinone.co.kr/ticker/?currency=" + strings.ToLower(currencyCode))

	if err != nil {
		log.Println("Error calling coinone api: ", err)
		return SearchResult{}, ErrSearchAPI
	}

	var data protocol.CoinoneTickerResponse
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Println("Error decoding coinone response: ", err)
		return SearchResult{}, ErrJSONParse
	}

	coinoneCache[currencyName] = SearchResult{
		LastPrice:    int(data.Last),
		Change:       int(data.Last - data.First),
		ChangeRate:   (data.Last - data.First) / data.First * 100,
		CurrencyName: currencyName,
		Time:         time.Now(),
	}
	return coinoneCache[currencyName], nil
}

/**
 * SearchAll: search for current currency exchange infomation for all supported currencies
 */
func (coinone) SearchAll() ([]SearchResult, error) {
	resultList := []SearchResult{}

	if time.Now().Sub(lastCoinoneCacheUpdate) < minTimeElapse {
		for currencyName := range coinoneCache {
			resultList = append(resultList, coinoneCache[currencyName])
		}
		return resultList, nil
	}

	netClient := &http.Client{
		Timeout: timeout,
	}
	response, err := netClient.Get("https://api.coinone.co.kr/ticker/?currency=all")
	if err != nil {
		log.Println("Error calling bithumb api: ", err)
		return []SearchResult{}, ErrSearchAPI
	}

	var coinoneResp map[string]interface{}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&coinoneResp); err != nil {
		log.Println("Error decoding coinone response: ", err)
		return []SearchResult{}, ErrJSONParse
	}

	for currencyName, currencyCode := range CurrencyCodes {
		currencyInfo := coinoneResp[strings.ToLower(currencyCode)].(map[string]interface{})

		lastPrice, _ := strconv.Atoi(currencyInfo["last"].(string))
		firstPrice, _ := strconv.Atoi(currencyInfo["first"].(string))

		coinoneCache[currencyName] = SearchResult{
			LastPrice:    lastPrice,
			Change:       lastPrice - firstPrice,
			ChangeRate:   float64(lastPrice-firstPrice) / float64(firstPrice) * 100,
			CurrencyName: currencyName,
			Time:         time.Now(),
		}

		lastCoinoneCacheUpdate = time.Now()
		resultList = append(resultList, coinoneCache[currencyName])
	}

	return resultList, nil
}
