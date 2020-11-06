package sync

import (
	"fmt"
	"os"
	"strings"
)

const (
	scheduleKeyPrefix = "SCHEDULE_"
	pagerDutyTokenKey = "PAGERDUTY_TOKEN"
	slackTokenKey     = "SLACK_TOKEN"
)

type Config struct {
	Schedules      []Schedule
	PagerDutyToken string
	SlackToken     string
}

type Schedule struct {
	ScheduleId             string
	AllOnCallGroupName     string
	CurrentOnCallGroupName string
}

func NewConfigFromEnv() (*Config, error) {

	config := &Config{
		PagerDutyToken: os.Getenv(pagerDutyTokenKey),
		SlackToken:     os.Getenv(slackTokenKey),
	}

	for _, key := range os.Environ() {
		if strings.HasPrefix(key, scheduleKeyPrefix) {
			value := strings.Split(key, "=")[1]
			scheduleValues := strings.Split(value, ",")
			if len(scheduleValues) != 2 {
				return nil, fmt.Errorf("expecting schedule value to be a comma separated scheduleId,name but got %s", value)
			}
			config.Schedules = append(config.Schedules, Schedule{
				ScheduleId:             scheduleValues[0],
				AllOnCallGroupName:     fmt.Sprintf("all-oncall-%ss", scheduleValues[1]),
				CurrentOnCallGroupName: fmt.Sprintf("current-oncall-%s", scheduleValues[1]),
			})
		}
	}

	if len(config.Schedules) == 0 {
		return nil, fmt.Errorf("expecting at least one schedule defined as an env var using prefix SCHEDULE_")
	}

	return config, nil
}
