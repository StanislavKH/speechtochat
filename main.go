package main

import (
	"context"
	"fmt"
	"log"

	stc "github.com/StanislavKH/speechtochat/pkg/stc"
)

func main() {
	// Replace with the path to your JSON key file
	keyFilePath := "keyfile.json"
	// Replace with your ChatGPT API key
	chatGPTAPIKey := ""

	ctx := context.Background()

	stcService, err := stc.NewStcService(ctx, keyFilePath, chatGPTAPIKey)
	if err != nil {
		panic(err)
	}

	transcript, err := stcService.TranscribeAudio(ctx, "sampleaudio/audiofile.mp3")
	if err != nil {
		log.Fatalf("Error transcribing audio: %v", err)
	}

	chatResponse, err := stcService.SendChatRequest(ctx, transcript)
	if err != nil {
		log.Fatalf("Error sending chat request: %v", err)
	}

	fmt.Println("ChatGPT Response:", chatResponse)
}
