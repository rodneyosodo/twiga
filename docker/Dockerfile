# Copyright (c) 2024 rodneyosodo
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
# http://www.apache.org/licenses/LICENSE-2.0

FROM golang:1.22-alpine AS builder
ARG SVC
ARG GOARCH
ARG GOARM
ARG VERSION
ARG TIME
ARG COMMIT
ARG COMMIT_DATE

WORKDIR /go/src/github.com/rodneyosodo/twiga
COPY . .
RUN apk update \
    && apk add make\
    && make $SVC \
    && mv build/$SVC /exe

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
