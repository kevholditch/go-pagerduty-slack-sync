package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	scheduleKeyPrefix        = "SCHEDULE_"
	pluralizeAllOnCallGroup  = "PLURALIZE_ALL_ONCALL_GROUP"
	currentOnCallGroupPrefix = "CURRENT_ALL_ONCALL_GROUP_PREFIX"
	pagerDutyTokenKey        = "PAGERDUTY_TOKEN"
	slackTokenKey            = "SLACK_TOKEN"
	runInterval              = "RUN_INTERVAL_SECONDS"
	pdScheduleLookaheadKey   = "PAGERDUTY_SCHEDULE_LOOKAHEAD"
	runIntervalDefault       = 60
)

// Config is used to configure application
// PagerDutyToken - token used to connect to pagerduty API
// SlackToken - token used to connect to Slack API
type Config struct {
	Schedules                  []Schedule
	PluralizeAllOnCallGroup    bool
	CurrentOnCallGroupPrefix   string
	PagerDutyToken             string
	SlackToken                 string
	RunIntervalInSeconds       int
	PagerdutyScheduleLookahead time.Duration
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

	pluralizeStr, ok := os.LookupEnv(pluralizeAllOnCallGroup)
	if !ok {
	    pluralizeStr = "true"
	}
	pluralize, err := strconv.ParseBool(pluralizeStr)
	if err != nil {
	    return nil, fmt.Errorf("failed to parse %s as bool: %w", pluralizeAllOnCallGroup, err)
	}
	config.PluralizeAllOnCallGroup = pluralize

	currentGroupPrefix, ok := os.LookupEnv(currentOnCallGroupPrefix)
	if !ok {
	    currentGroupPrefix = "current-"
	}
	config.CurrentOnCallGroupPrefix = currentGroupPrefix

	runInterval := os.Getenv(runInterval)
	v, err := strconv.Atoi(runInterval)
	if err == nil {
		config.RunIntervalInSeconds = v
	}

	pagerdutyScheduleLookahead, err := getPagerdutyScheduleLookahead()
	if err != nil {
		return nil, err
	}
	config.PagerdutyScheduleLookahead = pagerdutyScheduleLookahead

	for _, key := range os.Environ() {
		if strings.HasPrefix(key, scheduleKeyPrefix) {
			value := strings.Split(key, "=")[1]
			scheduleValues := strings.Split(value, ",")
			if len(scheduleValues) != 2 {
				return nil, fmt.Errorf("expecting schedule value to be a comma separated scheduleId,name but got %s", value)
			}

			schedule := Schedule{
				ScheduleID:             scheduleValues[0],
				AllOnCallGroupName:     fmt.Sprintf("all-oncall-%s", scheduleValues[1]),
				CurrentOnCallGroupName: fmt.Sprintf("oncall-%s", scheduleValues[1]),
			}

			if config.PluralizeAllOnCallGroup {
			    schedule.AllOnCallGroupName += "s"
			}

			if config.CurrentOnCallGroupPrefix != "" {
			    schedule.CurrentOnCallGroupName = config.CurrentOnCallGroupPrefix + schedule.CurrentOnCallGroupName
			}

			config.Schedules = append(config.Schedules, schedule)
		}
	}

	if len(config.Schedules) == 0 {
		return nil, fmt.Errorf("expecting at least one schedule defined as an env var using prefix SCHEDULE_")
	}

	return config, nil
}

func getPagerdutyScheduleLookahead() (time.Duration, error) {
	result := time.Hour * 24 * 100

	pdScheduleLookahead, ok := os.LookupEnv(pdScheduleLookaheadKey)
	if !ok {
		return result, nil
	}

	v, err := time.ParseDuration(pdScheduleLookahead)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as time.Duration: %w", pdScheduleLookahead, err)
	}

	return v, nil
}
