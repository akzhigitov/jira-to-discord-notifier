package handler

import (
	"bugReporter/model"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type DiscordHandler struct {
	webHookUrl string
}

func NewDiscordHandler(webHookUrl string) DiscordHandler {
	return DiscordHandler{
		webHookUrl: webHookUrl,
	}
}

func (handler DiscordHandler) SendMessage(message model.Message) error {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.Printf("Send %v", message)

	_, err = http.Post(handler.webHookUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
	return err
}
