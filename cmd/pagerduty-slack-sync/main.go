package main

import (
	"github.com/kevholditch/go-pagerduty-slack-sync/internal/sync"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	config, err := sync.NewConfigFromEnv()
	if err != nil {
		logrus.Errorf("could not parse config, error: %v", err)
		os.Exit(-1)
		return
	}

	err = sync.Schedules(config)
	if err != nil {
		logrus.Errorf("could not sync schedules, error: %v", err)
		os.Exit(-1)
		return
	}
}
