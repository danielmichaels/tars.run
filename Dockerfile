FROM danielmichaels/ci-toolkit AS toolkit
FROM debian:bookworm-slim AS litestream_downloader

ARG litestream_version="v0.3.13"
WORKDIR /litestream

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates \
      wget

RUN set -x && \
    if [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then \
      ARCH="arm7" ; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
      ARCH="arm64" ; \
    else \
      ARCH="amd64" ; \
    fi && \
    set -u && \
    litestream_binary_tgz_filename="litestream-${litestream_version}-linux-${ARCH}.tar.gz" && \
    wget "https://github.com/benbjohnson/litestream/releases/download/${litestream_version}/${litestream_binary_tgz_filename}" && \
    mv "${litestream_binary_tgz_filename}" litestream.tgz
RUN tar -xvzf litestream.tgz

FROM node:lts-slim AS node
COPY --from=toolkit ["/usr/local/bin/task", "/usr/local/bin/task"]

# PNPM is required to build the assets
RUN corepack enable pnpm
RUN mkdir -p /build
WORKDIR /build

COPY . .
RUN ["task", "assets"]

FROM golang:1.22-bookworm AS builder

WORKDIR /build
# only copy mod file for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY --from=node ["/build/assets/static", "/build/assets/static"]
COPY --from=toolkit ["/usr/local/bin/templ", "/usr/local/bin/templ"]
COPY --from=toolkit ["/usr/local/bin/task", "/usr/local/bin/task"]

COPY . .

#RUN ["task", "templgen"]

RUN apt-get install git -y &&\
    go build  \
    -ldflags="-s -w" \
    -o app ./cmd/app

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=toolkit ["/usr/local/bin/goose", "/usr/local/bin/goose"]
COPY --from=builder ["/build/entrypoint", "/app/entrypoint"]
COPY --from=builder ["/build/assets/migrations", "/app/migrations"]
COPY --from=builder ["/build/app", "/usr/bin/app"]
COPY --from=builder ["/build/litestream.yml", "/etc/litestream.yml"]
COPY --from=litestream_downloader ["/litestream/litestream", "/usr/bin/litestream"]

RUN mkdir -p data &&\
    apt-get update && apt-get install ca-certificates curl -y &&\
    chmod +x /app/entrypoint

ENTRYPOINT ["/app/entrypoint"]
