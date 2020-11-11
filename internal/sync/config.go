package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	scheduleKeyPrefix  = "SCHEDULE_"
	pagerDutyTokenKey  = "PAGERDUTY_TOKEN"
	slackTokenKey      = "SLACK_TOKEN"
	runInterval        = "RUN_INTERVAL_SECONDS"
	runIntervalDefault = 60
)

// Config is used to configure application
// PagerDutyToken - token used to connect to pagerduty API
// SlackToken - token used to connect to Slack API
type Config struct {
	Schedules            []Schedule
	PagerDutyToken       string
	SlackToken           string
	RunIntervalInSeconds int
}

// Schedule models a PagerDuty schedule that will be synced with Slack
// ScheduleID - PagerDuty schedule id to sync
// AllOnCallGroupName - Slack group name for all members of schedule
// CurrentOnCallGroupName - Slack group name for current person on call
type Schedule struct {
	ScheduleID             string
	AllOnCallGroupName     string
	CurrentOnCallGroupName string
}

// NewConfigFromEnv is a function to generate a config from env varibles
// PAGERDUTY_TOKEN - PagerDuty Token
// SLACK_TOKEN - Slack Token
// SCHEDULE_XXX="id,name" e.g. 1234,platform-engineer will generate a schedule with the following values
// ScheduleID = "1234", AllOnCallGroupName = "all-oncall-platform-engineers", CurrentOnCallGroupName: "current-oncall-platform-engineer"
func NewConfigFromEnv() (*Config, error) {

	config := &Config{
		PagerDutyToken:       os.Getenv(pagerDutyTokenKey),
		SlackToken:           os.Getenv(slackTokenKey),
		RunIntervalInSeconds: runIntervalDefault,
	}

	runInterval := os.Getenv(runInterval)
	v, err := strconv.Atoi(runInterval)
	if err == nil {
		config.RunIntervalInSeconds = v
	}

	for _, key := range os.Environ() {
		if strings.HasPrefix(key, scheduleKeyPrefix) {
			value := strings.Split(key, "=")[1]
			scheduleValues := strings.Split(value, ",")
			if len(scheduleValues) != 2 {
				return nil, fmt.Errorf("expecting schedule value to be a comma separated scheduleId,name but got %s", value)
			}
			config.Schedules = append(config.Schedules, Schedule{
				ScheduleID:             scheduleValues[0],
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
