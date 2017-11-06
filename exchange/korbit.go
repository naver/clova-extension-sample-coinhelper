package exchange

import (
	"coinHelper/protocol"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type korbit struct{}

// UseKorbit is a empty instance(?) of bithumb struct for using its methods
var (
	UseKorbit             korbit
	korbitCache           = map[string]SearchResult{}
	lastKorbitCacheUpdate time.Time
)

/**
 * Search: search for current currency exchange infomation for given currency
 */
func (korbit) Search(currencyName string) (SearchResult, error) {
	cachedResult := korbitCache[currencyName]

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
	response, err := netClient.Get("https://api.korbit.co.kr/v1/ticker/detailed?currency_pair=" + strings.ToLower(currencyCode) + "_krw")

	if err != nil {
		log.Println("Error calling bithumb api: ", err)
		return SearchResult{}, ErrSearchAPI
	}

	var data protocol.KorbitTickerData
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Println("Error decoding korbit response: ", err)
		return SearchResult{}, ErrJSONParse
	}

	korbitCache[currencyName] = SearchResult{
		LastPrice:    int(data.Last),
		Change:       int(data.Change),
		ChangeRate:   data.ChangePercent,
		CurrencyName: currencyName,
		Time:         time.Now(),
	}
	return korbitCache[currencyName], nil
}

/**
 * SearchAll: search for current currency exchange infomation for all supported currencies
 */
func (k korbit) SearchAll() ([]SearchResult, error) {
	resultList := []SearchResult{}
	if time.Now().Sub(lastKorbitCacheUpdate) < minTimeElapse {
		for currencyName := range korbitCache {
			resultList = append(resultList, korbitCache[currencyName])
		}
		return resultList, nil
	}

	prices := make(chan SearchResult)

	for currencyName := range CurrencyCodes {
		go func(currencyName string) {
			result, _ := k.Search(currencyName)
			prices <- result
		}(currencyName)
	}

	for {
		select {
		case result := <-prices:
			resultList = append(resultList, result)
			if len(resultList) == len(CurrencyCodes) {
				lastKorbitCacheUpdate = time.Now()
				return resultList, nil
			}
		}
	}

	return resultList, nil
}
