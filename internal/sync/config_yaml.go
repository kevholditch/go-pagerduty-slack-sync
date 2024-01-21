package sync

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type ConfigYaml struct {
	RunInterval                int           `yaml:"run-interval"`
	PagerDutyLookAheadDuration time.Duration `yaml:"pager-duty-look-ahead-duration"`
	PagerDutyToken             string        `yaml:"pager-duty-token"`
	SlackToken                 string        `yaml:"slack-token"`
	CurrentOnCallEnabled       bool          `yaml:"current-on-call-enabled"`
	CurrentOnCallNameFormat    string        `yaml:"current-on-call-name-format"`
	AllOnCallEnabled           bool          `yaml:"all-on-call-enabled"`
	AllOnCallNameFormat        string        `yaml:"all-on-call-name-format"`
	UserGroups                 []UserGroup   `yaml:"user-groups"`
}

type UserGroup struct {
	Name                    string   `yaml:"name"`
	ScheduleIDs             []string `yaml:"schedules"`
	CurrentOnCallEnabled    bool     `yaml:"current-on-call-enabled"`
	CurrentOnCallNameFormat string   `yaml:"current-on-call-name-format"`
	AllOnCallEnabled        bool     `yaml:"all-on-call-enabled"`
	AllOnCallNameFormat     string   `yaml:"all-on-call-name-format"`
}

func NewConfigFromYaml(filePath string) (*Config, error) {
	fmt.Printf("Reading YAML file %s\n", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var configYaml ConfigYaml
	err = yaml.Unmarshal(data, &configYaml)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	if len(configYaml.UserGroups) == 0 {
		return nil, fmt.Errorf("expecting at least one user-groups defined in YAML file")
	}

	defaultCurrentOnCallNameFormat := getOrDefault(configYaml.CurrentOnCallNameFormat, currentOnCallFormatDefault)
	defaultAllOnCallNameFormat := getOrDefault(configYaml.AllOnCallNameFormat, allOnCallFormatDefault)

	scheduleList := make([]Schedule, len(configYaml.UserGroups))

	fmt.Printf("Creating %d user-groups.\n", len(scheduleList))
	for i, s := range configYaml.UserGroups {
		fmt.Printf("Registering user-group %s.\n", s.Name)

		currentGroupName := fmt.Sprintf(getOrDefault(s.CurrentOnCallNameFormat, defaultCurrentOnCallNameFormat), s.Name)
		allGroupName := fmt.Sprintf(getOrDefault(s.AllOnCallNameFormat, defaultAllOnCallNameFormat), s.Name)

		fmt.Printf("Registering schedules(%q) for user-groups(%s, %s).\n", s.ScheduleIDs, currentGroupName, allGroupName)
		scheduleList[i] = Schedule{
			ScheduleIDs:            s.ScheduleIDs,
			CurrentOnCallGroupName: currentGroupName,
			AllOnCallGroupName:     allGroupName,
		}

	}

	runIntervalInSeconds := configYaml.RunInterval
	if runIntervalInSeconds == 0 {
		runIntervalInSeconds = runIntervalDefault
	}

	pagerdutyScheduleLookahead := configYaml.PagerDutyLookAheadDuration
	if pagerdutyScheduleLookahead == 0 {
		pagerdutyScheduleLookahead = pagerDutyLookAheadDefault
	}

	config := &Config{
		Schedules:                  scheduleList,
		PagerDutyToken:             configYaml.PagerDutyToken,
		SlackToken:                 configYaml.SlackToken,
		PagerdutyScheduleLookahead: pagerdutyScheduleLookahead,
		RunIntervalInSeconds:       runIntervalInSeconds,
		CurrentOnCallEnabled:       configYaml.CurrentOnCallEnabled,
		AllOnCallEnabled:           configYaml.AllOnCallEnabled,
	}

	return config, nil
}

func getOrDefault(format string, formatDefault string) string {
	if format == "" {
		return formatDefault
	}
	return format
}
