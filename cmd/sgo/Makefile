default: build install

.PHONY: build
build:
	go build -trimpath -ldflags='-extldflags -static -s -w' .

.PHONY: install
install:
	cp sgo $(HOME)/bin/sgo
	sudo setcap "CAP_DAC_OVERRIDE+eip" $(HOME)/bin/sgo
