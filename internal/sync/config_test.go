package sync

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfigFromEnv_SingleScheduleDefined(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform-engineer")()
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
	assert.Equal(t, "1234", config.Schedules[0].ScheduleID)
	assert.Equal(t, "all-oncall-platform-engineers", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "current-oncall-platform-engineer", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleID:             "1234",
		AllOnCallGroupName:     "all-oncall-platform-engineers",
		CurrentOnCallGroupName: "current-oncall-platform-engineer",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithCustomizedGroupNames(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform-team")()
	defer SetEnv("PLURALIZE_ALL_ONCALL_GROUP", "false")()
	defer SetEnv("CURRENT_ALL_ONCALL_GROUP_PREFIX", "")()
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
	assert.Equal(t, "1234", config.Schedules[0].ScheduleID)
	assert.Equal(t, "all-oncall-platform-team", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "oncall-platform-team", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleID:             "1234",
		AllOnCallGroupName:     "all-oncall-platform-team",
		CurrentOnCallGroupName: "oncall-platform-team",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithDefaultRunInterval(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform-engineer")()
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()

	config, err := NewConfigFromEnv()

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 60, config.RunIntervalInSeconds)
	assert.Equal(t, 1, len(config.Schedules))
	assert.Equal(t, "1234", config.Schedules[0].ScheduleID)
	assert.Equal(t, "all-oncall-platform-engineers", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "current-oncall-platform-engineer", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleID:             "1234",
		AllOnCallGroupName:     "all-oncall-platform-engineers",
		CurrentOnCallGroupName: "current-oncall-platform-engineer",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_SingleScheduleDefinedWithScheduleLookahead(t *testing.T) {
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform-engineer")()
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
	assert.Equal(t, "1234", config.Schedules[0].ScheduleID)
	assert.Equal(t, "all-oncall-platform-engineers", config.Schedules[0].AllOnCallGroupName)
	assert.Equal(t, "current-oncall-platform-engineer", config.Schedules[0].CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleID:             "1234",
		AllOnCallGroupName:     "all-oncall-platform-engineers",
		CurrentOnCallGroupName: "current-oncall-platform-engineer",
	}},
		config.Schedules))
}

func Test_NewConfigFromEnv_MultipleScheduleDefined(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_PLATFORM", "1234,platform-engineer")()
	defer SetEnv("SCHEDULE_CORE", "abcd,core-engineer")()
	defer SetEnv("SCHEDULE_UK", "efghass,uk-engineer")()

	config, err := NewConfigFromEnv()

	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(config.Schedules))

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleID:             "1234",
		AllOnCallGroupName:     "all-oncall-platform-engineers",
		CurrentOnCallGroupName: "current-oncall-platform-engineer",
	},
		{
			ScheduleID:             "abcd",
			AllOnCallGroupName:     "all-oncall-core-engineers",
			CurrentOnCallGroupName: "current-oncall-core-engineer",
		},
		{
			ScheduleID:             "efghass",
			AllOnCallGroupName:     "all-oncall-uk-engineers",
			CurrentOnCallGroupName: "current-oncall-uk-engineer",
		}},
		config.Schedules))
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
