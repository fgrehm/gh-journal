package ghjournal

import (
	"fmt"
	"github.com/octokit/go-octokit/octokit"
)

type GitHubClient interface {
	Events(pageNum int) ([]GitHubEvent, error)
}

type gitHubClient struct {
	userName string
	client   *octokit.Client
}

type GitHubEvent map[string]interface{}

func NewGitHubClient(userName, token string) GitHubClient {
	return &gitHubClient{
		userName: userName,
		client:   octokit.NewClient(octokit.TokenAuth{token}),
	}
}

func (c *gitHubClient) Events(pageNum int) ([]GitHubEvent, error) {
	url, err := octokit.UserURL.Expand(octokit.M{"user": c.userName})
	if err != nil {
		return nil, err
	}

	user, result := c.client.Users(url).One()
	if result.HasError() {
		return nil, err
	}

	eventsUrl, err := user.ReceivedEventsURL.Expand(nil)
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
