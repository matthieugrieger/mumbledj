all: mumbledj

mumbledj: main.go commands.go parseconfig.go strings.go song.go playlist.go songqueue.go
	go get github.com/nitrous-io/goop
	goop update
	goop go build .
		
clean:
	rm -f mumbledj
		
install:
	mkdir -p ~/.mumbledj/config
	mkdir -p ~/.mumbledj/songs
	if [ -a ~/.mumbledj/config/mumbledj.gcfg ]; then mv ~/.mumbledj/config/mumbledj.gcfg ~/.mumbledj/config/mumbledj_backup.gcfg; fi;
	cp -u mumbledj.gcfg ~/.mumbledj/config/mumbledj.gcfg
	sudo cp -f mumbledj /usr/local/bin/mumbledj

build:
	goop go build .
	
