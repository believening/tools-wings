default: build install

.PHONY: build
build:
	go build -ldflags='-extldflags -static -s -w' .
	sudo setcap "CAP_DAC_OVERRIDE+eip" sgo

.PHONY: install
install:
	cp sgo $(HOME)/bin/sgo
