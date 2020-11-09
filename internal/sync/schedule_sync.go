package sync

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"time"
	"github.com/slack-go/slack"
)

// Schedules does the sync
func Schedules(config *Config) error {

	p := pagerduty.NewClient(config.PagerDutyToken)
	users, err := p.ListOnCallUsers(config.Schedules[0].ScheduleID, pagerduty.ListOnCallUsersOptions{
		Since:         time.Now().UTC().Format(time.RFC3339),
		Until:         time.Now().UTC().Add(time.Hour*24*100).Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("%s %s\n", user.Name, user.Email)
	}

	s := slack.New(config.SlackToken)

	schedule := config.Schedules[0]
	g, err := s.CreateUserGroup(slack.UserGroup{
		Name:        schedule.CurrentOnCallGroupName,
	})

	if err != nil {
		return err
	}

	fmt.Printf(g.Name)

	return nil
}
