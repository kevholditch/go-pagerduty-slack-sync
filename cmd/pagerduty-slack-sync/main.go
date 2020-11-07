package main

import (
	"fmt"
	"github.com/kevholditch/go-pagerduty-slack-sync/internal/sync"
)

func main() {

	config, err := sync.NewConfigFromEnv()
	fmt.Print(config)
	fmt.Print(err)
}