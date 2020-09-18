package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	JiraUrl      = os.Getenv("JIRA_URL")
	WebHookUrl   = os.Getenv("WEB_HOOK_URL")
	JiraUsername = os.Getenv("JIRA_USERNAME")
	JiraPassword = os.Getenv("JIRA_PASSWORD")
	JiraFilterId = os.Getenv("JIRA_FILTER_ID")
)

type Message struct {
	Embeds []Embed `json:"embeds"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
	Author      Author `json:"author"`
}

type Author struct {
	Name string `json:"name"`
	Icon string `json:"icon_url"`
}

var (
	jiraClient = &jira.Client{}
)

func createMessageFromIssues(issues []jira.Issue) (Message, bool) {
	message := Message{}
	for _, issue := range issues {
		if issue.Fields.Watches.IsWatching {
			continue
		}

		embed := Embed{
			Title:       fmt.Sprintf("%s: %s", issue.Fields.Summary, issue.Key),
			Description: issue.Fields.Description,
			URL:         JiraUrl + "/browse/" + issue.Key,
			Color:       15746887, //red
			Author: Author{
				Name: issue.Fields.Creator.Name,
				Icon: issue.Fields.Priority.IconURL,
			},
		}

		message.Embeds = append(message.Embeds, embed)

		_, err := jiraClient.Issue.AddWatcher(issue.ID, JiraUsername)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return message, len(message.Embeds) > 0
}

func main() {
	tp := jira.BasicAuthTransport{
		Username: JiraUsername,
		Password: JiraPassword,
	}

	jiraClient, err := jira.NewClient(tp.Client(), JiraUrl)
	if err != nil {
		log.Fatal(err)
	}

	filterId, err := strconv.Atoi(JiraFilterId)
	if err != nil {
		log.Fatal(err)
	}
	filter, _, err := jiraClient.Filter.Get(filterId)
	if err != nil {
		log.Fatal(err)
	}
	issues, _, err := jiraClient.Issue.Search(filter.Jql, &jira.SearchOptions{})
	if err != nil {
		log.Fatal(err)
	}

	message, ok := createMessageFromIssues(issues)

	if ok {
		bytesRepresentation, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Send %v", string(bytesRepresentation))

		_, err = http.Post(WebHookUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
