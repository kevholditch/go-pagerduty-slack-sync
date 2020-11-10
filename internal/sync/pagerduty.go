package sync

import (
	"github.com/PagerDuty/go-pagerduty"
	"time"
)

type pagerDutyClient struct {
	client *pagerduty.Client
}

func newPagerDutyClient(token string) *pagerDutyClient {
	return &pagerDutyClient{
		client: pagerduty.NewClient(token),
	}
}

func (p *pagerDutyClient) getEmailsOfAllOnCallForSchedule(ID string) ([]string, error) {
	users, err := p.client.ListOnCallUsers(ID, pagerduty.ListOnCallUsersOptions{
		Since: time.Now().UTC().Format(time.RFC3339),
		Until: time.Now().UTC().Add(time.Hour * 24 * 100).Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	var results []string
	for _, user := range users {
		results = append(results, user.Email)
	}
	return results, nil
}

func (p *pagerDutyClient) getEmailsOfCurrentOnCallForSchedule(ID string) ([]string, error) {
	users, err := p.client.ListOnCallUsers(ID, pagerduty.ListOnCallUsersOptions{
		Since: time.Now().UTC().Format(time.RFC3339),
		Until: time.Now().UTC().Add(time.Second).Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	var results []string
	for _, user := range users {
		results = append(results, user.Email)
	}
	return results, nil
}
