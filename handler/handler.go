package handler

import (
	"coinHelper/intent"
	"coinHelper/protocol"
	"encoding/json"
	"log"
	"net/http"
)

const (
	msgError       = "서비스 연결이 원활하지 않습니다."
	msgUnavailable = "지원하지 않는 기능이에요."
	msgGreet       = "빗썸, 코인원, 코빗 거래소의 시세를 제공합니다."
	msgLeave       = ""
)

// ServeHTTP handles CEK requests
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var request protocol.CEKRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error during parsing Request JSON ")
		respondError(w, msgError)
		return
	}

	switch request.Request.Type {
	case "LaunchRequest":
		respondSuccess(w, protocol.MakeCEKResponse(nil, handleLaunchRequest()))

	case "SessionEndedRequest":
		respondSuccess(w, protocol.MakeCEKResponse(nil, handleEndRequest()))

	case "IntentRequest":
		intentName := request.Request.Intent.Name
		slots := request.Request.Intent.Slots

		if !request.Session.New {
			sessionIntentName, sessionSlots := parseSessionAttributes(request.Session.SessionAttributes)
			if len(sessionIntentName) > 0 {
				intentName = sessionIntentName
			}

			for key, value := range sessionSlots {
				slots[key] = value
			}
		}

		if response, err := handleIntent(intentName, slots, request.Session.New); err == nil {
			respondSuccess(w, response)
		} else {
			respondError(w, msgError)
		}

	default:
		respondError(w, msgUnavailable)
	}
}

func handleLaunchRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeOutputSpeech(msgGreet),
		ShouldEndSession: false,
	}
}

func handleEndRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeOutputSpeech(msgLeave),
		ShouldEndSession: true,
	}
}

func handleIntent(intentName string, slots map[string]protocol.CEKSlot, shouldEndSession bool) (protocol.CEKResponse, error) {
	switch intentName {
	case "AskCoinPriceIntent":
		response, err := intent.HandleAskCoinIntent(slots)
		if err == nil {
			response.Response.ShouldEndSession = shouldEndSession
		}
		return response, err
	default:
		return intent.UseGuideIntent.Handle(slots)
	}
	return protocol.CEKResponse{}, nil
}

func parseSessionAttributes(sessionAttributes map[string]string) (intent string, slots map[string]protocol.CEKSlot) {
	slots = map[string]protocol.CEKSlot{}

	for key, value := range sessionAttributes {
		if key == "intent" {
			intent = value
		} else {
			slots[key] = protocol.CEKSlot{
				Name:  key,
				Value: value,
			}
		}
	}

	return intent, slots
}

func respondError(w http.ResponseWriter, msg string) {
	response := protocol.MakeCEKResponse(nil,
		protocol.CEKResponsePayload{
			OutputSpeech:     protocol.MakeOutputSpeech(msg),
			ShouldEndSession: true,
		})

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func respondSuccess(w http.ResponseWriter, response protocol.CEKResponse) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {}
