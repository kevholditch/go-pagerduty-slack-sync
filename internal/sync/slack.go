package sync

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

type slackClient struct {
	users      []slack.User
	userGroups []slack.UserGroup
	client     *slack.Client
}

func newSlackClient(token string) (*slackClient, error) {
	s := slack.New(token)

	userGroups, err := s.GetUserGroups()
	if err != nil {
		return nil, err
	}

	users, err := s.GetUsers()
	if err != nil {
		return nil, err
	}

	return &slackClient{
		users:      users,
		userGroups: userGroups,
		client:     s,
	}, nil
}

func (s *slackClient) createOrGetUserGroup(name string) (*slack.UserGroup, error) {
	group := s.findUserGroupByName(name)
	if group != nil {
		return group, nil
	}

	g, err := s.client.CreateUserGroup(slack.UserGroup{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return &g, err
}

func (s *slackClient) getSlackIDsFromEmails(emails []string) ([]string, error) {
	var results []string
	for _, email := range emails {
		ID := s.findUserIDByEmail(email)
		if ID == nil {
			return nil, fmt.Errorf("could not find slack user with email: %s", email)
		}
		results = append(results, *ID)
	}
	return results, nil
}

func (s *slackClient) findUserIDByEmail(email string) *string {
	for _, u := range s.users {
		if strings.EqualFold(email, u.Profile.Email) {
			return &u.ID
		}
	}
	return nil
}

func (s *slackClient) findUserGroupByName(name string) *slack.UserGroup {
	for _, g := range s.userGroups {
		if strings.EqualFold(name, g.Name) {
			return &g
		}
	}
	return nil
}
