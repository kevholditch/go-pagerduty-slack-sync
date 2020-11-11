# go-pagerduty-slack-sync

This tool syncs one or more pagerduty schedules to two slack groups per schedule.  The first slack group is a group with everyone who is in the pagerduty schedule and the second contains the current person on call.

For example if you have the following users in a pagerduty schedule:

```
schedule id: 1234
user1
user2
user3 <= currently on call
```

Then when you run a sync for the schedule using:

```
docker run -e RUN_INTERVAL_SECONDS=60 -e SLACK_TOKEN=xxx -e PAGERDUTY_TOKEN=xxx -e SCHEDULE_PLATFORM=1234,platform-engineer kevholditch/pagerduty-slack-sync:latest
```

The following slack groups would be created:
    - @all-oncall-platform-engineers - user1, user2, user3 
    - @current-oncall-platform-engineer - user3
    
Multiple schedules can be synced at once by passing many env variables beginning with `SCHEDULE_`.  The format for the value that the schedule parameter expects is `<pagerduty schedule id>/<group-name>`.  The `<group name>` will be used to build the two names for the slack groups using the following format:
    - `all-oncall-<group-name>s`
    - `current-oncall-<group-name>`