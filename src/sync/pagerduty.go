package sync

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type pagerDutyClient struct {
	client *pagerduty.Client
}

func newPagerDutyClient(token string) *pagerDutyClient {
	return &pagerDutyClient{
		client: pagerduty.NewClient(token),
	}
}

func (p *pagerDutyClient) getEmailsForSchedule(ID string, lookahead time.Duration) ([]string, error) {
	users, err := p.client.ListOnCallUsers(ID, pagerduty.ListOnCallUsersOptions{
		Since: time.Now().UTC().Format(time.RFC3339),
		Until: time.Now().UTC().Add(lookahead).Format(time.RFC3339),
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
