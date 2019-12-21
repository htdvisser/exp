FROM alpine:3

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache ca-certificates

ADD dist/echo-server-$TARGETOS-$TARGETARCH /echo-server

CMD ["/echo-server"]
