package main

import (
	"bugReporter/handler"
	"bugReporter/utils"
	"log"
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
	f, err := os.OpenFile("bugReporter", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
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

	messages := jiraHandler.CreateMessageFromIssues(issues, utils.String2Map(LabelsRoles, " ", ":"))

	discordHandler := handler.NewDiscordHandler(WebHookUrl)
	for _, message := range messages {
		err := discordHandler.SendMessage(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
