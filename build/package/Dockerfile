ARG GOVERSION=1.15.4
FROM golang:${GOVERSION} as builder
COPY . /go-pagerduty-slack-sync
WORKDIR /go-pagerduty-slack-sync
ENV GO111MODULE=on CGO_ENABLED=0
RUN go build ./cmd/pagerduty-slack-sync/

FROM gcr.io/distroless/base
LABEL org.opencontainers.image.authors="<Kevin Holditch> kevholditch@hotmail.com"
WORKDIR /app
COPY --from=builder /go-pagerduty-slack-sync/pagerduty-slack-sync /app
CMD ["./pagerduty-slack-sync"]
