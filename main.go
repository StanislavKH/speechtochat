package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// Replace with the path to your JSON key file
	//keyFilePath := "keyfile.json"
	// Replace with your ChatGPT API key
	//chatGPTAPIKey := ""

	/*
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
	*/
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: true,
	}))
	e.GET("/ws", handleWebSocket)
	e.Start(":8080")
}

func handleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				return nil
			}
			log.Println("WebSocket read error:", err)
			break
		}
		// Process audio chunk here, e.g., save to a file or analyze
		if messageType == websocket.BinaryMessage {
			fileId := getUUID()
			log.Println("Received audio chunk:", len(p), "bytes")
			err = ioutil.WriteFile("sampleaudio/"+fileId+".wav", p, 0644)
			if err != nil {
				log.Println("Error writing audio file:", err)
				continue
			}
			log.Printf("Audio saved with ID: %s\n", fileId)
		}
	}

	return nil
}

func getUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}
