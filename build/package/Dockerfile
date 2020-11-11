FROM golang:1.15.4-buster as builder
COPY . /go-pagerduty-slack-sync
WORKDIR /go-pagerduty-slack-sync
ENV GO111MODULE=on
RUN GOOS=linux go build ./cmd/pagerduty-slack-sync/

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /go-pagerduty-slack-sync/pagerduty-slack-sync /app
CMD ["./pagerduty-slack-sync"]