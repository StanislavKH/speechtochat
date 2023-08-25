package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/StanislavKH/speechtochat/internal/stcsvc"
	"github.com/StanislavKH/speechtochat/pkg/stc"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var stcService *stc.StcService

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var err error
	keyFilePath := "keyfile.json"
	chatGPTAPIKey := ""

	ctx := context.Background()

	stcService, err = stc.NewStcService(ctx, keyFilePath, chatGPTAPIKey)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: true,
	}))
	e.GET("/get-records", getRecordsHandler)
	e.GET("/process-record/:elementName", processRecordHandler)
	e.DELETE("/remove-element/:elementName", removeRecordHandler)
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

func getRecordsHandler(c echo.Context) error {
	folderPath := "sampleaudio"
	fileInfoList, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice to store the file names
	fileNames := make([]string, 0)

	// Iterate through the list of file info and extract file names
	for _, fileInfo := range fileInfoList {
		if !fileInfo.IsDir() {
			fileNames = append(fileNames, fileInfo.Name())
		}
	}
	return c.JSON(http.StatusOK, fileNames)
}

func removeRecordHandler(c echo.Context) error {
	recordName := c.Param("elementName")
	folderPath := "./sampleaudio"

	filePath := filepath.Join(folderPath, recordName)

	err := os.Remove(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	response := map[string]interface{}{"message": "Record removed", "element": recordName}
	return c.JSON(http.StatusOK, response)
}

func processRecordHandler(c echo.Context) error {
	recordName := c.Param("elementName")

	stt, chatrsp, err := stcsvc.TranscribeSpeecToChat(*stcService, recordName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	response := map[string]interface{}{"stt": stt, "analysis": chatrsp, "element": recordName}
	return c.JSON(http.StatusOK, response)
}

func getUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}
