package intent

import (
	"coinHelper/exchange"
	"coinHelper/protocol"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	msgAskMarketName        = "어떤 거래소의 시세를 알려드릴까요?"
	msgAskCurrencyName      = "어떤 가상화폐 시세를 알려드릴까요?"
	msgAskAll               = "알고싶은 가상화폐 이름과 거래소 이름을 다시한번 정확히 말해주세요."
	msgMarketNotAvailable   = "말씀하신 거래소는 지원하지 않아요. 빗썸, 코빗, 코인원 중 한가지를 골라주세요."
	msgCurrencyNotAvailable = "%s에서는 %s 거래를 지원하지 않아요. 다른 가상화폐를 물어보세요."
	msgSearchResultUp       = "%s 거래소의 %s 시세입니다. 현재 %s원으로 24시간 기준 %s원, 약 %.2f%% 상승했습니다."
	msgSearchResultDown     = "%s 거래소의 %s 시세입니다. 현재 %s원으로 24시간 기준 %s원, 약 %.2f%% 하락했습니다."
	intentName              = "AskIntent"
)

// HandleAskCoinIntent handles intent request and response apropriate result
func HandleAskCoinIntent(slots map[string]protocol.CEKSlot) (response protocol.CEKResponse, err error) {
	marketName := slots["market"].Value
	currencyName := slots["currency"].Value
	if len(marketName) == 0 && len(currencyName) == 0 {
		marketName = getRandomMarket()
	}

	if len(currencyName) == 0 {
		response, err = getPriceFromMarket(marketName)
		return
	}

	if len(marketName) == 0 && len(currencyName) != 0 {
		response, err = getPricesByCurrencyType(currencyName)
		return
	}

	response, err = getCurrencyPriceFromMarket(currencyName, marketName)
	return
}

func getRandomMarket() string {
	var marketNames = []string{}
	for marketName := range exchangeMarkets {
		marketNames = append(marketNames, marketName)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	randomIndex := rand.Intn(len(marketNames))

	return marketNames[randomIndex]
}

func getPricesByCurrencyType(currencyName string) (response protocol.CEKResponse, err error) {
	type priceResults struct {
		Price      int
		MarketName string
	}

	prices := make(chan priceResults)
	priceList := []string{}

	for marketName, market := range exchangeMarkets {
		go func(marketName string, market exchange.Exchange) {
			result, _ := market.Search(currencyName)
			prices <- priceResults{
				Price:      result.LastPrice,
				MarketName: marketName,
			}
		}(marketName, market)
	}

	for {
		select {
		case price := <-prices:
			priceList = append(priceList, fmt.Sprintf("%s %s원", price.MarketName, humanize.Comma(int64(price.Price))))
			if len(priceList) == len(exchangeMarkets) {
				outputSpeechText := fmt.Sprintf("거래소별 %s 시세 입니다. %s 입니다.", currencyName, strings.Join(priceList, ", "))
				return makeResponse(nil, outputSpeechText, 0, true), nil
			}
		}
	}

	return
}

func getPriceFromMarket(marketName string) (response protocol.CEKResponse, err error) {
	var (
		outputSpeechText  string
		price             int
		sessionAttributes map[string]string
		shouldEndSession  = true
	)

	defer func() {
		response = makeResponse(sessionAttributes, outputSpeechText, price, shouldEndSession)
	}()

	market := exchangeMarkets[marketName]
	if market == nil {
		outputSpeechText = msgMarketNotAvailable
		return
	}

	priceList := []string{}
	result, _ := market.SearchAll()
	for _, data := range result {
		priceList = append(priceList, fmt.Sprintf("%s %s원", data.CurrencyName, humanize.Comma(int64(data.LastPrice))))
	}

	outputSpeechText = fmt.Sprintf("%s 거래소 시세 입니다.  %s 입니다.", marketName, strings.Join(priceList, ", "))
	response = makeResponse(sessionAttributes, outputSpeechText, price, shouldEndSession)
	return
}

func getCurrencyPriceFromMarket(currencyName string, marketName string) (response protocol.CEKResponse, err error) {
	var (
		outputSpeechText  string
		price             int
		sessionAttributes map[string]string
		shouldEndSession  = true
	)

	market := exchangeMarkets[marketName]
	if market == nil {
		outputSpeechText = msgMarketNotAvailable
		return
	}

	if len(exchange.CurrencyCodes[currencyName]) == 0 {
		outputSpeechText = fmt.Sprintf(msgCurrencyNotAvailable, marketName, currencyName)
		return
	}

	result, searchError := market.Search(currencyName)

	if searchError != nil {
		err = searchError
		return
	}

	price = result.LastPrice
	if result.Change >= 0 {
		outputSpeechText = fmt.Sprintf(msgSearchResultUp, marketName, currencyName, humanize.Comma(int64(price)), humanize.Comma(int64(result.Change)), result.ChangeRate)
	} else {
		outputSpeechText = fmt.Sprintf(msgSearchResultDown, marketName, currencyName, humanize.Comma(int64(price)), humanize.Comma(int64(result.Change)), result.ChangeRate)
	}

	response = makeResponse(sessionAttributes, outputSpeechText, price, shouldEndSession)
	return
}

func makeResponse(sessionAttributes map[string]string, outputSpeechText string, price int, shouldEndSession bool) protocol.CEKResponse {
	var card interface{}

	if price > 0 {
		card = protocol.MakePriceTextTemplate(price)
	}

	return protocol.MakeCEKResponse(
		sessionAttributes,
		protocol.CEKResponsePayload{
			OutputSpeech:     protocol.MakeOutputSpeech(outputSpeechText),
			Card:             card,
			ShouldEndSession: shouldEndSession,
		})
}
