all: mumbledj

mumbledj: main.go commands.go parseconfig.go strings.go
	go build .
		
clean:
	rm -f mumbledj
		
install:
	mkdir -p ~/.mumbledj/config
	mkdir -p ~/.mumbledj/songs
	if [ -a ~/.mumbledj/config/mumbledj.gcfg ]; then mv ~/.mumbledj/config/mumbledj.gcfg ~/.mumbledj/config/mumbledj_backup.gcfg; fi;
	cp -u mumbledj.gcfg ~/.mumbledj/config/mumbledj.gcfg
	sudo cp -f mumbledj /usr/local/bin/mumbledj
	
install_deps:
	go get -u github.com/layeh/gumble/gumble
	go get -u github.com/layeh/gumble/gumbleutil
	go get -u github.com/layeh/gumble/gumble_ffmpeg
	go get -u code.google.com/p/gcfg
