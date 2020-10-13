package handler

import (
	"../model"
	"fmt"
	"github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

const DescriptionLenMax = 2048

type JiraHandler struct {
	jiraClient   *jira.Client
	JiraUrl      string
	JiraUsername string
}

func NewJiraHandler(jiraUsername string, jiraPassword string, jiraUrl string) (*JiraHandler, error) {
	tp := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraPassword,
	}

	jiraClient, err := jira.NewClient(tp.Client(), jiraUrl)
	if err != nil {
		return nil, err
	}

	return &JiraHandler{
		jiraClient:   jiraClient,
		JiraUrl:      jiraUrl,
		JiraUsername: jiraUsername,
	}, nil
}

func (handler JiraHandler) IssuesFromFilter(JiraFilterId string) ([]jira.Issue, error) {

	filterId, err := strconv.Atoi(JiraFilterId)
	if err != nil {
		return nil, err
	}
	filter, _, err := handler.jiraClient.Filter.Get(filterId)
	if err != nil {
		return nil, err
	}
	issues, _, err := handler.jiraClient.Issue.Search(filter.Jql, &jira.SearchOptions{})

	return issues, err
}

func (handler JiraHandler) createMessageContent(roles []string) string {
	return strings.Join(roles, " ")
}

func (handler JiraHandler) CreateMessageFromIssues(issues []jira.Issue, labelsRoles map[string]string) []model.Message {
	embedMessagesByContent := map[string][]model.Embed{}
	for _, issue := range issues {
		if issue.Fields.Watches.IsWatching {
			continue
		}

		embed := handler.createEmbed(issue)
		roles := handler.getRoles(issue.Fields.Labels, labelsRoles)
		content := handler.createMessageContent(roles)

		embedMessagesByContent[content] = append(embedMessagesByContent[content], embed)

		handler.markAsWatched(issue)
	}

	messages := handler.createMessages(embedMessagesByContent)

	return messages
}

func (handler JiraHandler) markAsWatched(issue jira.Issue) {
	_, err := handler.jiraClient.Issue.AddWatcher(issue.ID, handler.JiraUsername)
	if err != nil {
		log.Fatalln(err)
	}
}

func (handler JiraHandler) createMessages(embedMessagesByContent map[string][]model.Embed) []model.Message {
	var result []model.Message
	for content, embeds := range embedMessagesByContent {
		message := model.Message{
			Embeds:  embeds,
			Content: content,
		}

		result = append(result, message)
	}
	return result
}

func (handler JiraHandler) getRoles(labels []string, labelsRoles map[string]string) []string {
	sort.Strings(labels)
	var roles []string
	for _, label := range labels {
		role, ok := labelsRoles[strings.ToLower(label)]
		if !ok {
			roles = append(roles, label)
		} else {
			roles = append(roles, fmt.Sprintf("<@&%v>", role))
		}
	}
	return roles
}

func (handler JiraHandler) createEmbed(issue jira.Issue) model.Embed {
	description := fmt.Sprintf(
		"**Приоритет: %s**\nОписание:\n%s", issue.Fields.Priority.Name, issue.Fields.Description)
	if len(description) > DescriptionLenMax {
		description = description[:DescriptionLenMax]
	}

	embed := model.Embed{
		Title:       fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary),
		Description: description,
		URL:         handler.JiraUrl + "/browse/" + issue.Key,
		Color:       15746887, //red
		Author: model.Author{
			Name: issue.Fields.Creator.Name,
		},
	}
	return embed
}
