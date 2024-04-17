package sync

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfigFromEnv_SingleScheduleDefined(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("RUN_INTERVAL_SECONDS", "10")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 10, config.RunIntervalInSeconds)
	assert.Equal(t, time.Hour*24*100, config.PagerdutyScheduleLookahead)
	assert.Equal(t, 1, len(config.Schedules))
	assert.Equal(t, "team-platform-support", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "platform-support", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs:            []string{"1234"},
		AllOnCallGroupName:     "team-platform-support",
		CurrentOnCallGroupName: "platform-support",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithDefaultRunInterval(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 60, config.RunIntervalInSeconds)
	assert.Equal(t, 1, len(config.Schedules))
	assert.Equal(t, "team-platform-support", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "platform-support", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs:            []string{"1234"},
		AllOnCallGroupName:     "team-platform-support",
		CurrentOnCallGroupName: "platform-support",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithScheduleLookahead(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("PAGERDUTY_SCHEDULE_LOOKAHEAD", "8760h")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 60, config.RunIntervalInSeconds)
	assert.Equal(t, time.Hour*24*365, config.PagerdutyScheduleLookahead)
	assert.Equal(t, 1, len(config.Schedules))
	assert.Equal(t, "team-platform-support", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "platform-support", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs:            []string{"1234"},
		AllOnCallGroupName:     "team-platform-support",
		CurrentOnCallGroupName: "platform-support",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_MultipleScheduleDefined(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform")()
	defer SetEnv("SCHEDULE_website", "abcd,website")()
	defer SetEnv("SCHEDULE_MOBILE", "efghass,mobile")()

	config, err := NewConfigFromEnv()

	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(config.Schedules))

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{
		{
			ScheduleIDs:            []string{"1234"},
			AllOnCallGroupName:     "team-platform-support",
			CurrentOnCallGroupName: "platform-support",
		},
		{
			ScheduleIDs:            []string{"abcd"},
			AllOnCallGroupName:     "team-website-support",
			CurrentOnCallGroupName: "website-support",
		},
		{
			ScheduleIDs:            []string{"efghass"},
			AllOnCallGroupName:     "team-mobile-support",
			CurrentOnCallGroupName: "mobile-support",
		},
	},
		config.Schedules))
}

func Test_NewConfigFromEnv_WithScheduleGroups(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_website_1", "aaaa,website")()
	defer SetEnv("SCHEDULE_website_2", "bbbb,website")()
	defer SetEnv("SCHEDULE_website_3", "cccc,website")()

	config, err := NewConfigFromEnv()

	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Schedules))

	assert.EqualValues(t, []Schedule{
		{
			ScheduleIDs:            []string{"aaaa", "bbbb", "cccc"},
			AllOnCallGroupName:     "team-website-support",
			CurrentOnCallGroupName: "website-support",
		},
	},
		config.Schedules)
}

func Test_NewConfigFromEnv_NoSchedulesDefined(t *testing.T) {
	config, err := NewConfigFromEnv()

	assert.Errorf(t, err, "expecting at least one schedule defined as an env var using prefix SCHEDULE_")
	assert.Nil(t, config)
}

func Test_NewConfigFromEnv_InvalidScheduleData(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "foo,bar,buzz")()

	config, err := NewConfigFromEnv()

	assert.Errorf(t, err, "expecting schedule value to be a comma separated scheduleId,name but got foo,bar,buzz")
	assert.Nil(t, config)
}

func SetEnv(key, value string) func() {
	_ = os.Setenv(key, value)
	return func() {
		_ = os.Unsetenv(key)
	}
}
