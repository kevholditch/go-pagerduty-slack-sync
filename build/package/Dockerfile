# Use the latest available Go 1.21 image
FROM golang:1.21 as builder
COPY . /go-pagerduty-slack-sync
WORKDIR /go-pagerduty-slack-sync
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build ./cmd/pagerduty-slack-sync/

FROM gcr.io/distroless/base
LABEL maintainer="Kevin Holditch <kevholditch@hotmail.com>"
WORKDIR /app
COPY --from=builder /go-pagerduty-slack-sync/pagerduty-slack-sync /app
CMD ["./pagerduty-slack-sync"]
