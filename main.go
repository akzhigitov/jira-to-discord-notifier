package main

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	JiraUsername string
	JiraPassword string

	jiraURL      = os.Getenv("JIRA_URL")
	webHookURL   = os.Getenv("WEB_HOOK_URL")
	jiraFilterID = os.Getenv("JIRA_FILTER_ID")
	labelsRoles  = os.Getenv("LABELS_ROLES")
	apiKey       = os.Getenv("API_KEY")
	schedule     = os.Getenv("SCHEDULE")
)

type reminder struct {
}

func (r reminder) Run() {
	jiraHandler, err := NewJiraHandler(JiraUsername, JiraPassword, jiraURL)
	if err != nil {
		log.Fatal(err)
	}

	issues, err := jiraHandler.IssuesFromFilter(jiraFilterID)
	if err != nil {
		log.Fatal(err)
	}

	messages := jiraHandler.CreateMessageFromIssues(issues, String2Map(labelsRoles, ";", ":"))

	log.Infoln("Messages count:", len(messages))

	discordHandler := NewDiscordHandler(webHookURL, apiKey)
	for _, message := range messages {
		err := discordHandler.SendMessage(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	routes := gin.Default()

	reminder := reminder{}
	jobrunner.Now(reminder)

	jobrunner.Start()
	err := jobrunner.Schedule(schedule, reminder)
	if err != nil {
		log.Fatal(err)
	}

	err = routes.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
