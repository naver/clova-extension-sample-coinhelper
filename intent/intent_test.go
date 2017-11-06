package intent

import (
	"coinHelper/protocol"
	"strings"
	"testing"
)

func getOutputSpeechmsg(response protocol.CEKResponse) string {
	return response.Response.OutputSpeech.Values.Value
}

func TestAskIntent(t *testing.T) {
	var testAskIntent askIntent

	var (
		slots    = map[string]protocol.CEKSlot{}
		msg      string
		response protocol.CEKResponse
	)
	response, _ = testAskIntent.Handle(slots)
	msg = getOutputSpeechmsg(response)
	if msg != msgAskAll {
		t.Errorf("Error: msg should be \"%s\"", msgAskAll)
	}

	market := &protocol.CEKSlot{
		Value: "빗썸",
		Name:  "market",
	}

	currency := &protocol.CEKSlot{
		Value: "비트코인",
		Name:  "currency",
	}

	slots["market"] = *market

	response, _ = testAskIntent.Handle(slots)
	msg = getOutputSpeechmsg(response)
	if msg != msgAskCurrencyName {
		t.Errorf("Error: msg should be \"%s\"", msgAskCurrencyName)
	}

	slots["currency"] = *currency

	response, _ = testAskIntent.Handle(slots)
	msg = getOutputSpeechmsg(response)
	if response, _ := testAskIntent.Handle(slots); !strings.Contains(getOutputSpeechmsg(response), market.Value) {
		t.Errorf("Error: msg should contain market name \"%s\"", getOutputSpeechmsg(response))
	}

}
