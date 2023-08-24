package main

/*var upgrader = websocket.Upgrader{
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

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.File("web/index.html")
	})
	e.GET("/ws", handleWebSocket)
	e.Start(":8080")
}

func handleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	fmt.Println("WebSocket connection established")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			break
		}
		if messageType == websocket.BinaryMessage {
			fmt.Println("Received audio chunk:", len(p), "bytes")
			// Process audio chunk here, e.g., save to a file or analyze
		}
	}

	return nil
}*/

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err, len(p), messageType)
			return
		}

		fmt.Printf("received %d bytes, type: %d ", len(p), messageType)
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe("0.0.0.0:8080", nil)
}
