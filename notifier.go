package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	JiraUrl      string `json:"jiraUrl"`
	FilterId     int    `json:"filterId"`
	WebHookUrl   string `json:"webHookUrl"`
	JiraUsername string `json:"jiraUsername"`
	JiraPassword string `json:"jiraPassword"`
}

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
	config     = &Config{}
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
			URL:         config.JiraUrl + "/browse/" + issue.Key,
			Color:       15746887, //red
			Author: Author{
				Name: issue.Fields.Creator.Name,
				Icon: issue.Fields.Priority.IconURL,
			},
		}

		message.Embeds = append(message.Embeds, embed)

		_, err := jiraClient.Issue.AddWatcher(issue.ID, config.JiraUsername)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return message, len(message.Embeds) > 0
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(dir, "config.json")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln("cant read config file:", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("cant parse config:", err)
	}

	tp := jira.BasicAuthTransport{
		Username: config.JiraUsername,
		Password: config.JiraPassword,
	}

	jiraClient, err := jira.NewClient(tp.Client(), config.JiraUrl)
	if err != nil {
		log.Fatal(err)
	}

	filter, _, err := jiraClient.Filter.Get(config.FilterId)
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

		_, err = http.Post(config.WebHookUrl, "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
