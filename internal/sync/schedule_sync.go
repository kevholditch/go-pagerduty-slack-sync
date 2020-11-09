package sync

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"strings"
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
	fmt.Printf("slack token: %s\n", config.SlackToken)

	schedule := config.Schedules[0]
	_, err = s.CreateUserGroup(slack.UserGroup{
		Name:        schedule.CurrentOnCallGroupName,
	})

	// ignore if group already exists
	if err != nil && err.Error() != "name_already_exists" {
		return err
	}

	members, err := s.GetUserGroupMembers(schedule.CurrentOnCallGroupName)

	fmt.Printf("%v+", members)

	g, err := s.GetUserGroups()
	if err != nil {
		return err
	}

	group, err := findUserGroup(schedule.CurrentOnCallGroupName, g)
	if err != nil {
		return err
	}

	_, err = s.UpdateUserGroupMembers(group.ID,"U53FWU333")
	if err != nil {
		return err
	}

	return nil
}

func findUserGroup(name string, groups []slack.UserGroup) (*slack.UserGroup, error) {

	for _, g := range groups {
		if strings.EqualFold(name, g.Name) {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("could not find group: %s", name)

}
