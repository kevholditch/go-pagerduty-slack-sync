# Use the latest available Go 1.21 image
FROM --platform=linux/amd64 golang:1.21 as builder
COPY . /pagerduty-slack-sync
WORKDIR /pagerduty-slack-sync
ENV GO111MODULE=on GOOS=linux GOARCH=amd64
RUN CGO_ENABLED=0 go build ./cmd/pagerduty-slack-sync/

FROM --platform=linux/amd64 gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /pagerduty-slack-sync/pagerduty-slack-sync /app
CMD ["./pagerduty-slack-sync"]
