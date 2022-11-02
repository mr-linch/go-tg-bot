#syntax=docker/dockerfile:1.2

# build static binary
FROM golang:1.19.3-alpine3.16 as builder

WORKDIR /go/src/go-tg-bot

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_VERSION
ARG BUILD_REF
ARG BUILD_TIME

# compile
RUN CGO_ENABLED=0 go build \
    -ldflags="-w -s -extldflags \"-static\" -X \"main.buildVersion=${BUILD_VERSION}\" -X \"main.buildCommit=${BUILD_REF}\" -X \"main.buildTime=${BUILD_TIME}\"" \
    -a \
    -tags timetzdata \
    -o /bin/go-tg-bot .


# run
FROM alpine:3.16


COPY --from=builder /bin/go-tg-bot /bin/go-tg-bot

EXPOSE 8000

# Reference: https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.source="https://github.com/mr-linch/go-tg-bot"

ENTRYPOINT [ "/bin/go-tg-bot" ]