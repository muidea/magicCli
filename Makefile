.PHONY: build install clean

APP_NAME=magicCli
SRC=application/magicCli/cmd/main.go

build:
	go build -o bin/$(APP_NAME) $(SRC)

install:
	cp bin/$(APP_NAME) ~/.local/bin/

clean:
	rm -rf bin/
