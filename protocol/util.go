package protocol

import (
	"github.com/dustin/go-humanize"
)

// MakeCEKResponse creates CEKResponse instance with given params
func MakeCEKResponse(sessionAttributes map[string]string, responsePayload CEKResponsePayload) CEKResponse {
	response := CEKResponse{
		SessionAttributes: sessionAttributes,
		Response:          responsePayload,
	}

	return response
}

// MakeOutputSpeech creates OutputSpeech instance with given params
func MakeOutputSpeech(msg string) OutputSpeech {
	return OutputSpeech{
		Type: "SimpleSpeech",
		Values: Value{
			Lang:  "ko",
			Value: msg,
			Type:  "PlainText",
		},
	}
}

// MakePriceTextTemplate creates OutputSpeech instance with given params
func MakePriceTextTemplate(price int) Card {
	var textCard = Card{}
	textCard.Type = "Text"
	textCard.HighlightText.Type = "number"
	textCard.HighlightText.Value = humanize.Comma(int64(price))
	textCard.MainText.Type = "string"
	textCard.MainText.Value = "원"

	return textCard
}
