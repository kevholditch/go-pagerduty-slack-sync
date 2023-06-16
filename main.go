package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/loveholidays/go-pagerduty-slack-sync/src/sync"
	"github.com/sirupsen/logrus"
)

func main() {

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	config, err := sync.NewConfigFromEnv()
	if err != nil {
		logrus.Errorf("could not parse config, error: %v", err)
		os.Exit(-1)
		return
	}

	logrus.Infof("starting, going to sync %d schedules", len(config.Schedules))

	timer := time.NewTicker(time.Second * time.Duration(config.RunIntervalInSeconds))

	for alive := true; alive; {
		select {
		case <-stop:
			logrus.Infof("stopping...")
			alive = false
			os.Exit(0)
		case <-timer.C:
			err = sync.Schedules(config)
			if err != nil {
				logrus.Errorf("could not sync schedules, error: %v", err)
				os.Exit(-1)
				return
			}
		}
	}

}
