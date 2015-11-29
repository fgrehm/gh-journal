package ghjournal

import (
	"fmt"

	"github.com/octokit/go-octokit/octokit"
)

type GitHubEvent map[string]interface{}

type GitHubClient interface {
	Events(pageNum int) ([]GitHubEvent, error)
}

type gitHubClient struct {
	user *octokit.User
	client   *octokit.Client
}

func NewGitHubClient(userName, token string) (GitHubClient, error) {
	client := octokit.NewClient(octokit.TokenAuth{token})

	url, err := octokit.UserURL.Expand(octokit.M{"user": userName})
	if err != nil {
		return nil, err
	}

	user, result := client.Users(url).One()
	if result.HasError() {
		return nil, err
	}

	return &gitHubClient{
		user:     user,
		client:   client,
	}, nil
}

func (c *gitHubClient) Events(pageNum int) ([]GitHubEvent, error) {
	eventsUrl, err := c.user.ReceivedEventsURL.Expand(nil)
	if err != nil {
		return nil, err
	}

	// REFACTOR: There should be a better way of building an URL per page
	receivedEventsUrl := fmt.Sprintf("%s?page=%d", eventsUrl.String(), pageNum)
	req, err := c.client.NewRequest(receivedEventsUrl)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	_, err = req.Get(&data)
	if err != nil {
		return nil, err
	}
	events := []GitHubEvent{}
	for _, event := range data {
		events = append(events, event)
	}

	return events, nil
}
