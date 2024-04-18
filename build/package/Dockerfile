# Use the latest available Go 1.21 image
FROM golang:1.21 as builder
COPY . /pagerduty-slack-sync
WORKDIR /pagerduty-slack-sync
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build ./cmd/pagerduty-slack-sync/

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /pagerduty-slack-sync/pagerduty-slack-sync /app
CMD ["./pagerduty-slack-sync"]
