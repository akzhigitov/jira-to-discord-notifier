package main

import (
	"fmt"
	"github.com/maraticus/jira-to-discord-notifier/handler"
	"github.com/maraticus/jira-to-discord-notifier/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	jiraURL      = os.Getenv("JIRA_URL")
	webHookURL   = os.Getenv("WEB_HOOK_URL")
	jiraUsername = os.Getenv("JIRA_USERNAME")
	jiraPassword = os.Getenv("JIRA_PASSWORD")
	jiraFilterID = os.Getenv("JIRA_FILTER_ID")
	labelsRoles  = os.Getenv("LABELS_ROLES")
)

func main() {
	f, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	jiraHandler, err := handler.NewJiraHandler(jiraUsername, jiraPassword, jiraURL)
	if err != nil {
		log.Fatal(err)
	}

	issues, err := jiraHandler.IssuesFromFilter(jiraFilterID)
	if err != nil {
		log.Fatal(err)
	}

	messages := jiraHandler.CreateMessageFromIssues(issues, utils.String2Map(labelsRoles, ";", ":"))

	discordHandler := handler.NewDiscordHandler(webHookURL)
	for _, message := range messages {
		err := discordHandler.SendMessage(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
