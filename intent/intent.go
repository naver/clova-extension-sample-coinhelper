package intent

import (
	"coinHelper/exchange"
	"coinHelper/protocol"
)

type intent interface {
	Handle(slots map[string]protocol.CEKSlot) (string, error)
}

var exchangeMarkets = map[string]exchange.Exchange{
	"코빗":  exchange.UseKorbit,
	"빗썸":  exchange.UseBithumb,
	"코인원": exchange.UseCoinone,
}
