package exchange

import (
	"errors"
	"time"
)

// Exchange is an interface for crypto currency exchange market
type Exchange interface {
	Search(string) (SearchResult, error)
	SearchAll() ([]SearchResult, error)
}

// ErrCoinTypeNotAvailable is an error instance for invalid currency type
var (
	ErrCoinTypeNotAvailable = errors.New("Error: Requested currency type is not available for this exchange market")
	ErrSearchAPI            = errors.New("Error: Search API returned an error or API server is not available")
	ErrJSONParse            = errors.New("Error: An error occured while trying to parse JSON response from Search API")
	CurrencyCodes           = map[string]string{
		"비트코인":     "BTC",
		"이더리움":     "ETH",
		"이더리움 클래식": "ETC",
		"리플":       "XRP",
		"비트코인 캐시":  "BCH",
	}
)

// SearchResult contains API call result
type SearchResult struct {
	LastPrice    int
	Change       int
	ChangeRate   float64
	CurrencyName string
	Time         time.Time
}

const minTimeElapse = 1 * time.Second
const timeout = time.Second * 2
