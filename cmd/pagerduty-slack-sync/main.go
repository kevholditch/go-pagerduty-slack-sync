package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bitrise-io/pagerduty-slack-sync/internal/sync"
	"github.com/sirupsen/logrus"
)

func main() {

	// Make the stop channel buffered
	stop := make(chan os.Signal, 1) // Buffer size of 1
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	config, err := sync.NewConfigFromEnv()
	if err != nil {
		logrus.Errorf("could not parse config, error: %v", err)
		os.Exit(-1)
		return
	}

	// Start the health check endpoint and make sure not to block
	go func() {
		_ = http.ListenAndServe(":8080", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
			},
		))
	}()

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
