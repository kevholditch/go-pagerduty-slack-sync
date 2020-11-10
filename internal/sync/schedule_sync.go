package sync

import (
	"fmt"
)

// Schedules does the sync
func Schedules(config *Config) error {

	s, err := newSlackClient(config.SlackToken)
	if err != nil {
		return err
	}
	p := newPagerDutyClient(config.PagerDutyToken)

	for _, schedule := range config.Schedules {

		emails, err := p.getEmailsOfOnCallForSchedule(schedule.ScheduleID)
		if err != nil {
			return err
		}
		slackIDs, err := s.getSlackIDsFromEmails(emails)
		if err != nil {
			return err
		}
		for _, ID := range slackIDs {
			fmt.Printf("%s\n", ID)
		}

	}

	//group, err := s.createOrGetUserGroup(schedule.CurrentOnCallGroupName)
	//if err != nil {
	//	return err
	//}
	//
	//members, err := s.GetUserGroupMembers(schedule.CurrentOnCallGroupName)
	//
	//fmt.Printf("%v+", members)
	//
	//_, err = s.UpdateUserGroupMembers(groupID,"U53FWU333")
	//if err != nil {
	//	return err
	//}

	return nil
}
