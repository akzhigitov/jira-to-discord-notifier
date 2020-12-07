package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DiscordHandler struct {
	webHookURL string
	apiKey     string
}

func NewDiscordHandler(webHookUrl string, apiKey string) DiscordHandler {
	return DiscordHandler{
		webHookURL: webHookUrl,
		apiKey:     apiKey,
	}
}

func (handler DiscordHandler) SendMessage(message Message) error {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.Infoln("Discord request", message)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", handler.webHookURL, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	response, err := client.Do(req)
	log.Infoln("Discord response", response)
	return err
}
