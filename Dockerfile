FROM golang:1.15.6-alpine as builder
LABEL stage=builder
WORKDIR /go/src/github.com/lexycore/gitlab-go-import
COPY . .
RUN apk add --no-cache --virtual .build-deps bash libarchive-tools git openssh-client ca-certificates openssl libffi-dev openssl-dev gcc musl-dev libc-dev make

RUN export GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 && \
    go mod download && \
    go mod vendor && \
    go generate -v ./... && \
    go install -tags netgo -ldflags '-w -extldflags "-static"' -v ./cmd/...

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/gitlab-go-import /
ENTRYPOINT ["/gitlab-go-import"]
