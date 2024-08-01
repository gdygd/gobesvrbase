  
.PHONY: fmt build test bench
.EXPORT_ALL_VARIABLES:

#BUILD_DATE = `date +%Y-%m-%d\ %H:%M`
#VERSIONFILE = ./version.go


#GO111MODULE ?= on

all: build hello

deps:
	go get github.com/ghetzel/shmtool/shm
	go get gopkg.in/ini.v1

test:
	go test -v

build:
	go build -o ../bin/mngprc mngprc.go

run:
	go run mngprc.go

