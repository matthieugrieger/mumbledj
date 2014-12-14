all: mumbledj

mumbledj: main.go commands.go parseconfig.go strings.go
		go build .
		
clean:
		rm -f mumbledj
		
install:
		sudo cp -f mumbledj /usr/local/bin/mumbledj
		mkdir -p ~/.mumbledj/config
		mkdir -p ~/.mumbledj/songs
		-mv ~/.mumbledj/config/config.toml ~/.mumbledj/config/config_backup.toml
		cp -u config.toml ~/.mumbledj/config/config.toml
			
install_deps:
		go get -u github.com/layeh/gumble/gumble
		go get -u github.com/layeh/gumble/gumbleutil
		go get -u github.com/layeh/gumble/gumble_ffmpeg
		go get -u github.com/BurntSushi/toml
