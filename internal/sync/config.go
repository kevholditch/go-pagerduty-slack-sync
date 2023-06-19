package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Env var mapping
const (
	scheduleKeyPrefix         = "SCHEDULE_"
	pagerDutyTokenKey         = "PAGERDUTY_TOKEN"
	slackTokenKey             = "SLACK_TOKEN"
	runInterval               = "RUN_INTERVAL_SECONDS"
	pdScheduleLookaheadKey    = "PAGERDUTY_SCHEDULE_LOOKAHEAD"
	oncallGroupNamePrefixKey  = "ON_CALL_GROUP_NAME_PREFIX"
	requiresAllOnCallGroupKey = "REQUIRES_ALL_ON_CALL_GROUP"
)

// Config default values
const (
	runIntervalDefault              = 60
	defaultCurrentOnCallGroupPrefix = "current-oncall"
)

// Group identifier enum mapping
const (
	allGroupTag     = "all"
	currentGroupTag = "current"
	noPluralTag     = "noPlural"
)

// Group Config
// Required - flag to denote if group should be created / synced
// Name - group name
type GroupConfig struct {
	IsRequired bool
	Name       string
}

// Schedule models a PagerDuty schedule that will be synced with Slack
// ScheduleIDs - All PagerDuty schedule ID's to sync
// AllOnCallGroupName - Slack group name for all members of schedule
// CurrentOnCallGroupName - Slack group name for current person on call
type Schedule struct {
	ScheduleIDs        []string
	AllOnCallGroup     GroupConfig
	CurrentOnCallGroup GroupConfig
}

// Config is used to configure application
// PagerDutyToken - token used to connect to pagerduty API
// SlackToken - token used to connect to Slack API
type Config struct {
	Schedules                  []Schedule
	PagerDutyToken             string
	SlackToken                 string
	RunIntervalInSeconds       int
	PagerdutyScheduleLookahead time.Duration
	OnCallGroupNamePrefix      string
}

// NewConfigFromEnv is a function to generate a config from env variables
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

	pagerdutyScheduleLookahead, err := getPagerdutyScheduleLookahead()
	if err != nil {
		return nil, err
	}
	config.PagerdutyScheduleLookahead = pagerdutyScheduleLookahead

	config.OnCallGroupNamePrefix = getOnCallGroupNamePrefix()

	for _, key := range os.Environ() {
		if strings.HasPrefix(key, scheduleKeyPrefix) {
			value := strings.Split(key, "=")[1]
			scheduleValues := strings.Split(value, ",")
			if len(scheduleValues) < 2 {
				return nil, fmt.Errorf("expecting schedule value to be a comma separated scheduleId,name but got %s", value)
			}

			requiresAllGroup, requiresCurrentGroup, err := determineGroupsRequired(scheduleValues)
			if err != nil {
				return nil, err
			}

			pluralize := determineGroupPluralization(scheduleValues)

			config.Schedules = appendSchedule(
				config.Schedules,
				config.OnCallGroupNamePrefix,
				scheduleValues[0],
				scheduleValues[1],
				requiresCurrentGroup,
				requiresAllGroup,
				pluralize,
			)
		}
	}

	if len(config.Schedules) == 0 {
		return nil, fmt.Errorf("expecting at least one schedule defined as an env var using prefix SCHEDULE_")
	}

	return config, nil
}

func determineGroupsRequired(scheduleValues []string) (bool, bool, error) {
	// all and current groups will be created by default if not specified
	requiresAllGroup := true
	requiresCurrentGroup := true

	hasGroupsSpecified := len(scheduleValues) >= 3 && len(strings.TrimSpace(scheduleValues[2])) > 0

	if hasGroupsSpecified {
		groups := scheduleValues[2]

		// if groups are specified by the user they will be created only if exactly specified
		requiresAllGroup = false
		requiresCurrentGroup = false

		groupNames := strings.Split(groups, "|")
		for _, group := range groupNames {
			switch group {
			case allGroupTag:
				if requiresAllGroup {
					return false, false, fmt.Errorf("Duplicate group identifier found %s", groups)
				}
				requiresAllGroup = true
			case currentGroupTag:
				if requiresCurrentGroup {
					return false, false, fmt.Errorf("Duplicate group identifier found %s", groups)
				}
				requiresCurrentGroup = true
			default:
				return false, false, fmt.Errorf("Unknown group identifier found, expected either %s or %s but got %s", allGroupTag, currentGroupTag, group)
			}
		}
	}

	return requiresAllGroup, requiresCurrentGroup, nil
}

func determineGroupPluralization(scheduleValues []string) bool {
	// all group name will be pluralized by default
	pluralize := true

	if len(scheduleValues) >= 4 && scheduleValues[3] == noPluralTag {
		pluralize = false
	}

	return pluralize
}

func appendSchedule(
	schedules []Schedule,
	prefix string,
	scheduleID string,
	teamName string,
	requiresCurrentGroup bool,
	requiresAllGroup bool,
	pluralize bool,
) []Schedule {
	allOnCallGroupSuffix := ""
	if pluralize {
		allOnCallGroupSuffix = "s"
	}

	currentGroupConfig := GroupConfig{
		IsRequired: requiresCurrentGroup,
		Name:       fmt.Sprintf("%s-%s", prefix, teamName),
	}

	allGroupConfig := GroupConfig{
		IsRequired: requiresAllGroup,
		Name:       fmt.Sprintf("all-oncall-%s%s", teamName, allOnCallGroupSuffix),
	}

	updated := false
	newScheduleList := make([]Schedule, len(schedules))
	for i, s := range schedules {
		if s.CurrentOnCallGroup.Name != currentGroupConfig.Name {
			newScheduleList[i] = s

			continue
		}

		updated = true

		newScheduleList[i] = Schedule{
			ScheduleIDs:        append(s.ScheduleIDs, scheduleID),
			CurrentOnCallGroup: currentGroupConfig,
			AllOnCallGroup:     allGroupConfig,
		}
	}

	if !updated {
		newScheduleList = append(newScheduleList, Schedule{
			ScheduleIDs:        []string{scheduleID},
			CurrentOnCallGroup: currentGroupConfig,
			AllOnCallGroup:     allGroupConfig,
		})
	}

	return newScheduleList
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

func getOnCallGroupNamePrefix() string {
	oncallGroupNamePrefix, ok := os.LookupEnv(oncallGroupNamePrefixKey)

	if !ok {
		logrus.Infof("%s not provided - defaulting to %s", oncallGroupNamePrefixKey, defaultCurrentOnCallGroupPrefix)
		return defaultCurrentOnCallGroupPrefix
	}

	return oncallGroupNamePrefix
}
