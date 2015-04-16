all: mumbledj

mumbledj: main.go commands.go parseconfig.go strings.go service.go service_youtube.go songqueue.go cache.go
	go get github.com/nitrous-io/goop
	rm -rf Goopfile.lock
	goop install
	goop go build

clean:
	rm -f mumbledj*

install:
	mkdir -p ~/.mumbledj/config
	mkdir -p ~/.mumbledj/songs
	if [ -a ~/.mumbledj/config/mumbledj.gcfg ]; then mv ~/.mumbledj/config/mumbledj.gcfg ~/.mumbledj/config/mumbledj_backup.gcfg; fi;
	cp -u config.gcfg ~/.mumbledj/config/mumbledj.gcfg
	if [ -d ~/bin ]; then cp -f mumbledj* ~/bin/mumbledj; else sudo cp -f mumbledj* /usr/local/bin/mumbledj; fi;

build:
	goop go build
