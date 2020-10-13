package main

import (
	"bugReporter/handler"
	"bugReporter/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	JiraUrl      = os.Getenv("JIRA_URL")
	WebHookUrl   = os.Getenv("WEB_HOOK_URL")
	JiraUsername = os.Getenv("JIRA_USERNAME")
	JiraPassword = os.Getenv("JIRA_PASSWORD")
	JiraFilterId = os.Getenv("JIRA_FILTER_ID")
	LabelsRoles  = os.Getenv("LABELS_ROLES")
)

func main() {
	f, err := os.OpenFile("testlogrus.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	jiraHandler, err := handler.NewJiraHandler(JiraUsername, JiraPassword, JiraUrl)
	if err != nil {
		log.Fatal(err)
	}

	issues, err := jiraHandler.IssuesFromFilter(JiraFilterId)
	if err != nil {
		log.Fatal(err)
	}

	messages := jiraHandler.CreateMessageFromIssues(issues, utils.String2Map(LabelsRoles, ";", ":"))

	discordHandler := handler.NewDiscordHandler(WebHookUrl)
	for _, message := range messages {
		err := discordHandler.SendMessage(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
