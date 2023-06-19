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
	assert.Equal(t, "all-oncall-platform-engineers", config.Schedules[0].AllOnCallGroup.Name)
	assert.Equal(t, "current-oncall-platform-engineer", config.Schedules[0].CurrentOnCallGroup.Name)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs: []string{"1234"},
		AllOnCallGroup: GroupConfig{
			IsRequired: true,
			Name:       "all-oncall-platform-engineers",
		},
		CurrentOnCallGroup: GroupConfig{
			IsRequired: true,
			Name:       "current-oncall-platform-engineer",
		},
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
	assert.Equal(t, "all-oncall-platform-engineers", config.Schedules[0].AllOnCallGroup.Name)
	assert.Equal(t, "current-oncall-platform-engineer", config.Schedules[0].CurrentOnCallGroup.Name)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs: []string{"1234"},
		AllOnCallGroup: GroupConfig{
			IsRequired: true,
			Name:       "all-oncall-platform-engineers",
		},
		CurrentOnCallGroup: GroupConfig{
			IsRequired: true,
			Name:       "current-oncall-platform-engineer",
		},
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
	assert.Equal(t, "all-oncall-platform-engineers", config.Schedules[0].AllOnCallGroup.Name)
	assert.Equal(t, "current-oncall-platform-engineer", config.Schedules[0].CurrentOnCallGroup.Name)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{{
		ScheduleIDs: []string{"1234"},
		AllOnCallGroup: GroupConfig{
			IsRequired: true,
			Name:       "all-oncall-platform-engineers",
		},
		CurrentOnCallGroup: GroupConfig{
			IsRequired: true,
			Name:       "current-oncall-platform-engineer",
		},
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

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{
		{
			ScheduleIDs: []string{"1234"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-platform-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-platform-engineer",
			},
		},
		{
			ScheduleIDs: []string{"abcd"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-core-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-core-engineer",
			},
		},
		{
			ScheduleIDs: []string{"efghass"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-uk-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-uk-engineer",
			},
		},
	},
		config.Schedules))
}

func Test_NewConfigFromEnv_WithScheduleGroups(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_CORE_1", "aaaa,core-engineer")()
	defer SetEnv("SCHEDULE_CORE_2", "bbbb,core-engineer")()
	defer SetEnv("SCHEDULE_CORE_3", "cccc,core-engineer")()

	config, err := NewConfigFromEnv()

	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(config.Schedules))

	assert.EqualValues(t, []Schedule{
		{
			ScheduleIDs: []string{"aaaa", "bbbb", "cccc"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-core-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-core-engineer",
			},
		},
	},
		config.Schedules)
}

func Test_NewConfigFromEnv_WithSpecifiedGroups(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_CORE_1", "aaaa,core-engineer,current")()
	defer SetEnv("SCHEDULE_CORE_2", "bbbb,platform-engineer,all")()
	defer SetEnv("SCHEDULE_CORE_3", "cccc,uk-engineer,current|all")()

	config, err := NewConfigFromEnv()

	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(config.Schedules))

	assert.EqualValues(t, []Schedule{
		{
			ScheduleIDs: []string{"aaaa"},
			AllOnCallGroup: GroupConfig{
				IsRequired: false,
				Name:       "all-oncall-core-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-core-engineer",
			},
		},
		{
			ScheduleIDs: []string{"bbbb"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-platform-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: false,
				Name:       "current-oncall-platform-engineer",
			},
		},
		{
			ScheduleIDs: []string{"cccc"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-uk-engineers",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-uk-engineer",
			},
		},
	},
		config.Schedules)
}

func Test_NewConfigFromEnv_WithDisablePluralization(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_CORE_1", "aaaa,core-engineer,all,noPlural")()
	defer SetEnv("SCHEDULE_CORE_2", "bbbb,platform-engineer,current|all,noPlural")()
	defer SetEnv("SCHEDULE_CORE_3", "cccc,uk-engineer,,noPlural")()
	defer SetEnv("SCHEDULE_CORE_4", "dddd,frontend-engineer,         ,noPlural")()

	config, err := NewConfigFromEnv()

	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(config.Schedules))

	assert.EqualValues(t, []Schedule{
		{
			ScheduleIDs: []string{"aaaa"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-core-engineer",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: false,
				Name:       "current-oncall-core-engineer",
			},
		},
		{
			ScheduleIDs: []string{"bbbb"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-platform-engineer",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-platform-engineer",
			},
		},
		{
			ScheduleIDs: []string{"cccc"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-uk-engineer",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-uk-engineer",
			},
		},
		{
			ScheduleIDs: []string{"dddd"},
			AllOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "all-oncall-frontend-engineer",
			},
			CurrentOnCallGroup: GroupConfig{
				IsRequired: true,
				Name:       "current-oncall-frontend-engineer",
			},
		},
	},
		config.Schedules)
}

func Test_NewConfigFromEnv_WithInvalidGroupName(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_CORE_1", "aaaa,core-engineer,invalid_group_name")()

	_, err := NewConfigFromEnv()

	assert.Errorf(t, err, "expected error stating unknown group identifier found")
}

func Test_NewConfigFromEnv_WithDuplicateGroupName(t *testing.T) {
	defer SetEnv("PAGERDUTY_TOKEN", "token1")()
	defer SetEnv("SLACK_TOKEN", "secretToken1")()
	defer SetEnv("SCHEDULE_CORE_1", "aaaa,core-engineer,all|all|current|current")()

	_, err := NewConfigFromEnv()

	assert.Errorf(t, err, "expected error stating duplicate group identifier found")
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
