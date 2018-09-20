FROM golang:1.10 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/reeganexe/redis-query
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /redis-query .

FROM alpine
COPY --from=builder /redis-query /

ENTRYPOINT ["/redis-query"]
