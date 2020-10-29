package main

import (
	"github.com/akzhigitov/jira-to-discord-notifier/handler"
	"github.com/akzhigitov/jira-to-discord-notifier/utils"
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
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
	schedule     = os.Getenv("SCHEDULE")
)

type reminder struct {
}

func (r reminder) Run() {
	jiraHandler, err := handler.NewJiraHandler(jiraUsername, jiraPassword, jiraURL)
	if err != nil {
		log.Fatal(err)
	}

	issues, err := jiraHandler.IssuesFromFilter(jiraFilterID)
	if err != nil {
		log.Fatal(err)
	}

	messages := jiraHandler.CreateMessageFromIssues(issues, utils.String2Map(labelsRoles, ";", ":"))

	log.Debugln("Messages count:", len(messages))

	discordHandler := handler.NewDiscordHandler(webHookURL)
	for _, message := range messages {
		err := discordHandler.SendMessage(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	routes := gin.Default()

	reminder:=reminder{}
	jobrunner.Now(reminder)

	jobrunner.Start()
	err := jobrunner.Schedule(schedule, reminder)
	if err != nil {
		log.Fatal(err)
	}

	routes.Run(":8080")
}
