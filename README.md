# go-pagerduty-slack-sync
A small app to sync pagerduty on call groups to slack groups


```
docker run -e RUN_INTERVAL_SECONDS=60 -e SLACK_TOKEN=xxx -e PAGERDUTY_TOKEN=xxx -e SCHEDULE_PLATFORM=1234,platform-engineer kevholditch/pagerduty-slack-sync:latest
```
