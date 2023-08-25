package stcsvc

import (
	"github.com/StanislavKH/speechtochat/pkg/stc"
)

func TranscribeSpeecToChat(stcService stc.StcService, fileName string) (string, string, error) {
	transcript, err := stcService.TranscribeAudio("sampleaudio/" + fileName)
	if err != nil {
		return "", "", err
	}

	chatResponse, err := stcService.SendChatRequest(transcript)
	if err != nil {
		return "", "", err
	}

	return transcript, "ChatGPT Response: " + chatResponse, nil
}
