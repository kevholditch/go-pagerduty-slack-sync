package main

import (
	"fmt"
	"github.com/kevholditch/go-pagerduty-slack-sync/internal/sync"
	"os"
)

func main() {

	config, err := sync.NewConfigFromEnv()
	if err != nil {
		fmt.Printf("could not parse config, error: %v", err)
		os.Exit(-1)
		return
	}

	err = sync.Schedules(config)
	if err != nil {
		fmt.Printf("could not sync schdules, error: %v", err)
		os.Exit(-1)
		return
	}
}
