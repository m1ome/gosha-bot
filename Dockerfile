FROM golang:latest as builder

ADD . /go/src/github.com/m1ome/gosha-bot

RUN set -x \
    && cd /go/src/github.com/m1ome/gosha-bot/ \
    && export VERSION=$(git rev-parse --verify HEAD) \
    && export LDFLAGS="-w -s -X main.Version=${VERSION}" \
    && CGO_ENABLED=0 go build -v -ldflags "${LDFLAGS}" -o /go/bin/gosha-bot

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/gosha-bot /gosha-bot
COPY --from=builder /go/src/github.com/m1ome/gosha-bot/assets /assets

WORKDIR /

CMD ["/gosha-bot", "--help"]