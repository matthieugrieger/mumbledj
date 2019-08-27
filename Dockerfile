# Build environment for mumbledj - golang alpine container
FROM golang:1.12-alpine3.9 AS builder
ARG branch=master

ENV GO111MODULE=on

RUN apk add --update ca-certificates make git build-base opus-dev
RUN git clone -b $branch --single-branch https://github.com/Reikion/mumbledj.git $GOPATH/src/go.reik.pl/mumbledj

# add assets, which will be bundled with binary
WORKDIR $GOPATH/src/go.reik.pl/mumbledj
COPY assets assets
RUN make && make install


# Export binary only from builder environment
FROM alpine:latest
RUN apk add --update ffmpeg openssl aria2 python3 && \
    # youtube-dl use /usr/bin/env python so we need to create symlink
    ln -s /usr/bin/python3 /usr/bin/python && \
    wget https://yt-dl.org/downloads/latest/youtube-dl -O /bin/youtube-dl && \
    chmod a+x /bin/youtube-dl
COPY --from=builder /usr/local/bin/mumbledj /usr/local/bin/mumbledj

# Drop to user level privileges
RUN addgroup -S mumbledj && adduser -S mumbledj -G mumbledj && chmod 750 /home/mumbledj
WORKDIR /home/mumbledj
USER mumbledj
RUN mkdir -p .config/mumbledj && \
    mkdir -p .cache/mumbledj

ENTRYPOINT ["/usr/local/bin/mumbledj"]
