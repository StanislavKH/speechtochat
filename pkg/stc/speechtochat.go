package stc

import (
	"context"
	"io/ioutil"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

type StcService struct {
	SpeechClient *speech.Client
	ChatClient   *openai.Client
}

func NewStcService(ctx context.Context, speechKeyFilePath, openAIKey string) (*StcService, error) {
	speechClient, err := speech.NewClient(ctx, option.WithCredentialsFile(speechKeyFilePath))
	if err != nil {
		return nil, err
	}
	chatClient := openai.NewClient(openAIKey)
	return &StcService{
		SpeechClient: speechClient,
		ChatClient:   chatClient,
	}, nil
}

func (stc *StcService) TranscribeAudio(ctx context.Context, audioFilePath string) (string, error) {
	// Read the audio file
	audioData, err := ioutil.ReadFile(audioFilePath)
	if err != nil {
		return "", err
	}

	// Configure the request
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        *speechpb.RecognitionConfig_ENCODING_UNSPECIFIED.Enum(),
			SampleRateHertz: 31000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	}

	// Perform the recognition
	resp, err := stc.SpeechClient.Recognize(ctx, req)
	if err != nil {
		return "", err
	}

	// Get the transcribed text
	if len(resp.Results) > 0 && len(resp.Results[0].Alternatives) > 0 {
		return resp.Results[0].Alternatives[0].Transcript, nil
	}

	return "", nil
}

func (stc *StcService) SendChatRequest(ctx context.Context, transcript string) (string, error) {
	resp, err := stc.ChatClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Please transcribe this: " + transcript,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
