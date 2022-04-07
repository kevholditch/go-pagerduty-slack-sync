package sync

import (
	"strings"
	"time"

	"github.com/kevholditch/go-pagerduty-slack-sync/internal/compare"
	"github.com/sirupsen/logrus"
)

// Schedules does the sync
func Schedules(config *Config) error {
	logrus.Infof("running schedule sync...")
	s, err := newSlackClient(config.SlackToken)
	if err != nil {
		return err
	}
	p := newPagerDutyClient(config.PagerDutyToken)

	updateSlackGroup := func(emails []string, groupName string) error {
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

	getEmailsForSchedules := func(schedules []string, lookahead time.Duration) ([]string, error) {
		var emails []string

		for _, sid := range schedules {
			e, err := p.getEmailsForSchedule(sid, lookahead)
			if err != nil {
				return nil, err
			}

			emails = appendIfMissing(emails, e...)
		}

		return emails, nil
	}

	for _, schedule := range config.Schedules {
		logrus.Infof("checking slack group: %s", schedule.CurrentOnCallGroupName)

		currentOncallEngineerEmails, err := getEmailsForSchedules(schedule.ScheduleIDs, time.Second)
		if err != nil {
			return err
		}

		err = updateSlackGroup(currentOncallEngineerEmails, schedule.CurrentOnCallGroupName)
		if err != nil {
			return err
		}

		logrus.Infof("checking slack group: %s", schedule.AllOnCallGroupName)

		allOncallEngineerEmails, err := getEmailsForSchedules(schedule.ScheduleIDs, config.PagerdutyScheduleLookahead)
		if err != nil {
			return err
		}

		err = updateSlackGroup(allOncallEngineerEmails, schedule.AllOnCallGroupName)
		if err != nil {
			return err
		}
	}

	return nil
}

func appendIfMissing(slice []string, items ...string) []string {
out:
	for _, i := range items {
		for _, ele := range slice {
			if ele == i {
				continue out
			}
		}
		slice = append(slice, i)
	}

	return slice
}
