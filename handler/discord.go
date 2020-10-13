package handler

import (
	"../model"
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
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

	log.Infoln("Discord request", message)
	response, err := http.Post(handler.webHookUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
	log.Infoln("Discord response", response)
	return err
}
