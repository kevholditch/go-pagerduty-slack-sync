package sync

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	filePath = "example.yaml"
)

func Test_NewConfigFromYaml_SingleScheduleDefined(t *testing.T) {
	config, err := NewConfigFromYaml(filePath)

	assert.NoError(t, err)
	assert.Equal(t, "token1", config.PagerDutyToken)
	assert.Equal(t, "secretToken1", config.SlackToken)
	assert.Equal(t, 600, config.RunIntervalInSeconds)
	assert.Equal(t, time.Hour*24, config.PagerdutyScheduleLookahead)
	assert.Equal(t, 2, len(config.Schedules))
	assert.Equal(t, true, config.CurrentOnCallEnabled)
	assert.Equal(t, false, config.AllOnCallEnabled)

	firstSchedule := config.Schedules[0]
	assert.Equal(t, "all-oncall-platform-engineers", firstSchedule.AllOnCallGroupName)
	assert.Equal(t, "current-oncall-platform-engineer", firstSchedule.CurrentOnCallGroupName)

	assert.True(t, assert.ObjectsAreEqualValues([]Schedule{
		{
			ScheduleIDs:            []string{"1234"},
			AllOnCallGroupName:     "all-oncall-platform-engineers",
			CurrentOnCallGroupName: "current-oncall-platform-engineer",
		}, {
			ScheduleIDs:            []string{"123", "12345"},
			AllOnCallGroupName:     "all-company-wide",
			CurrentOnCallGroupName: "pd-company-wide",
		}},
		config.Schedules))
}
