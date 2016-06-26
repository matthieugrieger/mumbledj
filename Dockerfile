FROM alpine:3.3

ENV GOPATH=/

RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.3/community" >> /etc/apk/repositories
RUN apk add --update ca-certificates go ffmpeg mercurial make opus-dev openal-soft-dev git build-base python
RUN apk upgrade

RUN wget http://yt-dl.org/downloads/latest/youtube-dl -O /bin/youtube-dl && chmod a+x /bin/youtube-dl

COPY . /src/github.com/matthieugrieger/mumbledj
COPY config.yaml /root/.config/mumbledj/config.yaml

WORKDIR /src/github.com/matthieugrieger/mumbledj

RUN make
RUN make install
RUN apk del go mercurial git build-base make && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/mumbledj"]
