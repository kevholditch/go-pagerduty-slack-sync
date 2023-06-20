<p align="center">
  <img src="./synclogo.png">
</p>



This tool syncs one or more pagerduty schedules to two slack groups per schedule.  The first slack group is a group with everyone who is in the pagerduty schedule and the second contains the current person on call.

For example if you have the following users in a pagerduty schedule:

```bash
schedule id: 1234
user1
user2
user3 <= currently on call
```

Then when you run a sync for the schedule using:

```bash
docker run -e RUN_INTERVAL_SECONDS=60 -e SLACK_TOKEN=xxx -e PAGERDUTY_TOKEN=xxx -e SCHEDULE_PLATFORM=1234,platform-engineer kevholditch/pagerduty-slack-sync:latest
```

The following slack groups would be created:

- `@all-oncall-platform-engineers` => `user1, user2, user3`
- `@current-oncall-platform-engineer` => `user3`

Multiple schedules can be synced at once by passing many env variables beginning with `SCHEDULE_`.  The format for the value that the schedule parameter expects is `<pagerduty schedule id>/<group-name>`.  The `<group name>` will be used to build the two names for the slack groups using the following format:
    - `all-oncall-<group-name>s`
    - `current-oncall-<group-name>`

If there are multiple schedules with the same `<group-name>` are defined, then slack groups contains the combined list of all people for all the given schedules.

For instance given following environment variables:

```bash
-e SCHEDULE_TEAM_1=abcd,platform-engineer -e SCHEDULE_TEAM_2=efgh,platform-engineer
```

This will result in a pair of slack groups with the combined users:

- `@all-oncall-platform-engineers` => combined list of all users in `abcd` and `efgh` schedules
- `@current-oncall-platform-engineer` => combined list of current on call users in `abcd` and `efgh` schedules

## You may also add addition optional positional specifiers to `SCHEDULE_` environment variables in positions 3 and 4:

The the parameters are as followed

```bash
  3: current|all: which will specify if the current or/and all groups should be created
    - These group names can only be either `current` or `all` and must be joined with the `|` symbol 
    - if not specified both `current` and `all` groups will be created
    - if any other value is passed here an error will be returned when processed
    - duplicate groups will result in an error being returned when processed

  4: noPlural: which will disable pluralization of the group name for the `all` group if specified 
    - if not specified the all group name will be pluralized
```

for example:

```bash
-e SCHEDULE_TEAM_1=abcd,platform-engineer,current|all
-e SCHEDULE_TEAM_2=efgh,core-team,all,noPlural
-e SCHEDULE_TEAM_3=ijkl,data-engineering,current
-e SCHEDULE_TEAM_3=mnop,frontend-engineer,,noPlural
```

Please note to omit a optional positional argument please something like the following
```bash
-e SCHEDULE_TEAM_3=mnop,frontend-engineer,,noPlural
```

will create:

- `@all-oncall-platform-engineer` and `@current-oncall-platform-engineers` for `SCHEDULE_TEAM_1`
- `@all-oncall-core-team` with `core-team` not pluralized only for  `SCHEDULE_TEAM_2`
- `@current-oncall-infrastructure` only for `SCHEDULE_TEAM_2`
- `@all-oncall-frontend-engineer` and `@current-oncall-frontend-engineer` with `frontend-engineer` not pluralized for `SCHEDULE_TEAM_1`

Full parameter list:

| Env Name                        | Description                                                                       | Default Value  | Example                 |
|:--------------------------------|:----------------------------------------------------------------------------------|:---------------|:------------------------|
| PAGERDUTY_TOKEN                 | Token used to talk to the PagerDuty API                                           | n/a            | xxxxx                   |
| SLACK_TOKEN                     | Token used to talk to Slack API                                                   | n/a            | xoxp-xxxxxx             |
| SCHEDULE_<NAME>                 | A PagerDuty schedule that you want to sync                                        | n/a            | 1234,platform-engineer  |
| RUN_INTERVAL_SECONDS            | Run a sync every X seconds                                                        | 60             | 300                     |
| PAGERDUTY_SCHEDULE_LOOKAHEAD    | How far into the future to evaluate Pagerduty schedules (Go time duration format) | 2400h          | 8760h                   |
| ON_CALL_GROUP_NAME_PREFIX       | Prefix to be used for the current on call slack group                             | current-oncall | oncall                  |

## Slack permissions

In order for the app to run you will need to create a bot with the following permissions:
```
usergroups:read
usergroups:write
users:read
users:read.email
```

If you have locked down your slack so only the admins can create groups then you have two options.  You can either create the slack groups up front and the app will use those or you can give the bot user auth and give it admin perssions:
```
admin.usergroups:read
admin.usergroups:write
usergroups:read
usergroups:write
users:read
users:read.email
```
