FROM golang:1.20-buster as builder
COPY . /go-pagerduty-slack-sync
WORKDIR /go-pagerduty-slack-sync
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build ./cmd/pagerduty-slack-sync/

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /go-pagerduty-slack-sync/pagerduty-slack-sync /app
CMD ["./pagerduty-slack-sync"]