package sync

import (
	"github.com/kevholditch/go-pagerduty-slack-sync/internal/compare"
	"github.com/sirupsen/logrus"
	"strings"
)

// Schedules does the sync
func Schedules(config *Config) error {

	s, err := newSlackClient(config.SlackToken)
	if err != nil {
		return err
	}
	p := newPagerDutyClient(config.PagerDutyToken)

	updateSchedule := func(emails []string, groupName string) error {
		slackIDs, err := s.getSlackIDsFromEmails(emails)
		if err != nil {
			return err
		}

		userGroup, err := s.createOrGetUserGroup(groupName)
		if err != nil {
			return err
		}
		members, err := s.Client.GetUserGroupMembers(userGroup.ID)
		if err != nil {
			return err
		}

		if !compare.Array(slackIDs, members) {
			logrus.Infof("member list %s needs updating...", groupName)
			_, err = s.Client.UpdateUserGroupMembers(userGroup.ID, strings.Join(slackIDs, ","))
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, schedule := range config.Schedules {
		logrus.Infof("checking slack group: %s", schedule.CurrentOnCallGroupName)
		emails, err := p.getEmailsOfCurrentOnCallForSchedule(schedule.ScheduleID)
		if err != nil {
			return err
		}

		err = updateSchedule(emails, schedule.CurrentOnCallGroupName)
		if err != nil {
			return err
		}

		logrus.Infof("checking slack group: %s", schedule.AllOnCallGroupName)
		emails, err = p.getEmailsOfAllOnCallForSchedule(schedule.ScheduleID)
		if err != nil {
			return err
		}

		err = updateSchedule(emails, schedule.AllOnCallGroupName)
		if err != nil {
			return err
		}

	}

	return nil
}
