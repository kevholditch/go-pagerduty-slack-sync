package sync

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	scheduleKeyPrefix      = "SCHEDULE_"
	pagerDutyTokenKey      = "PAGERDUTY_TOKEN"
	slackTokenKey          = "SLACK_TOKEN"
	runInterval            = "RUN_INTERVAL_SECONDS"
	pdScheduleLookaheadKey = "PAGERDUTY_SCHEDULE_LOOKAHEAD"
	currentOnCallFormat    = "USER_GROUP_CURRENT_ON_CALL_FORMAT"
	currentOnCallEnabled   = "USER_GROUP_CURRENT_ON_CALL_ENABLED"
	allOnCallFormat        = "USER_GROUP_ALL_ON_CALL_FORMAT"
	allOnCallEnabled       = "USER_GROUP_ALL_ON_CALL_ENABLED"

	pagerDutyLookAheadDefault  = time.Hour * 24 * 100
	runIntervalDefault         = 60
	currentOnCallFormatDefault = "current-oncall-%s"
	allOnCallFormatDefault     = "all-oncall-%ss"
)

// Config is used to configure application
// PagerDutyToken - token used to connect to pagerduty API
// SlackToken - token used to connect to Slack API
type Config struct {
	Schedules                  []Schedule
	PagerDutyToken             string
	SlackToken                 string
	RunIntervalInSeconds       int
	PagerdutyScheduleLookahead time.Duration
	CurrentOnCallEnabled       bool
	AllOnCallEnabled           bool
}

// Schedule models a PagerDuty schedule that will be synced with Slack
// ScheduleIDs - All PagerDuty schedule ID's to sync
// AllOnCallGroupName - Slack group name for all members of schedule
// CurrentOnCallGroupName - Slack group name for current person on call
type Schedule struct {
	ScheduleIDs            []string
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
		RunIntervalInSeconds: GetEnvInt(runInterval, runIntervalDefault),
		CurrentOnCallEnabled: GetEnvBool(currentOnCallEnabled, true),
		AllOnCallEnabled:     GetEnvBool(allOnCallEnabled, true),
	}

	pagerdutyScheduleLookahead, err := GetPagerdutyScheduleLookahead()
	if err != nil {
		return nil, err
	}

	config.PagerdutyScheduleLookahead = pagerdutyScheduleLookahead

	currentOnCallNameFormat := GetEnvStr(currentOnCallFormat, currentOnCallFormatDefault)
	allOnCallNameFormat := GetEnvStr(allOnCallFormat, allOnCallFormatDefault)

	for _, key := range os.Environ() {
		if strings.HasPrefix(key, scheduleKeyPrefix) {
			value := strings.Split(key, "=")[1]
			scheduleValues := strings.Split(value, ",")
			if len(scheduleValues) != 2 {
				return nil, fmt.Errorf("expecting schedule value to be a comma separated scheduleId,name but got %s", value)
			}

			config.Schedules = appendSchedule(config, scheduleValues[0], scheduleValues[1], currentOnCallNameFormat, allOnCallNameFormat)
		}
	}

	if len(config.Schedules) == 0 {
		return nil, fmt.Errorf("expecting at least one schedule defined as an env var using prefix SCHEDULE_")
	}

	return config, nil
}

func appendSchedule(config *Config, scheduleID, teamName string, currentOnCallNameFormat string, allOnCallNameFormat string) []Schedule {

	currentGroupName := fmt.Sprintf(currentOnCallNameFormat, teamName)
	allGroupName := fmt.Sprintf(allOnCallNameFormat, teamName)

	newScheduleList := make([]Schedule, len(config.Schedules))
	updated := false

	for i, s := range config.Schedules {
		if s.CurrentOnCallGroupName != currentGroupName {
			newScheduleList[i] = s

			continue
		}

		updated = true

		newScheduleList[i] = Schedule{
			ScheduleIDs:            append(s.ScheduleIDs, scheduleID),
			AllOnCallGroupName:     allGroupName,
			CurrentOnCallGroupName: currentGroupName,
		}
	}

	if !updated {
		newScheduleList = append(newScheduleList, Schedule{
			ScheduleIDs:            []string{scheduleID},
			AllOnCallGroupName:     allGroupName,
			CurrentOnCallGroupName: currentGroupName,
		})
	}

	return newScheduleList
}

func GetPagerdutyScheduleLookahead() (time.Duration, error) {

	pdScheduleLookahead, ok := os.LookupEnv(pdScheduleLookaheadKey)
	if !ok {
		return pagerDutyLookAheadDefault, nil
	}

	v, err := time.ParseDuration(pdScheduleLookahead)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as time.Duration: %w", pdScheduleLookahead, err)
	}

	return v, nil
}
