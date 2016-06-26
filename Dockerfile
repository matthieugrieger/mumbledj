FROM alpine:3.3

ENV GOPATH=/

RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.3/community" >> /etc/apk/repositories
RUN apk add --update ca-certificates go ffmpeg make build-base opus-dev python
RUN apk upgrade

RUN wget https://yt-dl.org/downloads/latest/youtube-dl -O /bin/youtube-dl && chmod a+x /bin/youtube-dl

COPY . /src/github.com/matthieugrieger/mumbledj
COPY config.yaml /root/.config/mumbledj/config.yaml

WORKDIR /src/github.com/matthieugrieger/mumbledj

RUN make
RUN make install
RUN apk del go make build-base && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/mumbledj"]
