FROM golang:1.12-alpine3.9

ENV GO111MODULE=on

RUN apk add --update ca-certificates ffmpeg make git build-base opus-dev python aria2 openssl
RUN apk upgrade

RUN wget https://yt-dl.org/downloads/latest/youtube-dl -O /bin/youtube-dl && chmod a+x /bin/youtube-dl

RUN git clone https://github.com/Reikion/mumbledj.git $GOPATH/src/go.reik.pl/mumbledj
WORKDIR $GOPATH/src/go.reik.pl/mumbledj

RUN make
RUN make install
COPY asset/config.yaml /root/.config/mumbledj/config.yaml

RUN apk del go make build-base && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/mumbledj"]
