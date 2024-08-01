  
.PHONY: fmt build test bench
.EXPORT_ALL_VARIABLES:

#GO111MODULE ?= on

all: build test

deps:
	go get github.com/ghetzel/shmtool/shm
	go get gopkg.in/ini.v1
	go get github.com/gdygd/goshm
	go get github.com/gdygd/goglib
	
test:
	go test -v

build:
	go build -o ../bin/cli cli.go

run:
	go run cli.go