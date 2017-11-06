package intent

import (
	"coinHelper/protocol"
)

type guideIntent struct{}

// UseGuideIntent is used for handling intents
var UseGuideIntent guideIntent

func (guideIntent) Handle(slots map[string]protocol.CEKSlot) (protocol.CEKResponse, error) {

	return protocol.MakeCEKResponse(
		nil,
		protocol.CEKResponsePayload{
			OutputSpeech:     protocol.MakeOutputSpeech("이 익스텐션은 가상화폐 거래소에서 현재 가상화폐 시세를 조회할 수 있습니다. 빗썸에서 비트코인 시세 알려줘 등으로 조회해보세요."),
			ShouldEndSession: true,
		}), nil
}
