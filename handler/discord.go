package handler

import (
	"bytes"
	"encoding/json"
	"github.com/akzhigitov/jira-to-discord-notifier/model"
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

	log.Debugln("Discord request", message)
	response, err := http.Post(handler.webHookUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
	log.Debugln("Discord response", response)
	return err
}
