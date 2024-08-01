  
.PHONY: fmt build test bench
.EXPORT_ALL_VARIABLES:

#GO111MODULE ?= on

all: build test

deps:
	go get gopkg.in/ini.v1
	go get github.com/gorilla/handlers
	go get github.com/gorilla/mux
	go get github.com/go-sql-driver/mysql
	go get github.com/godror/godror
	go get github.com/lib/pq
	go get github.com/fogleman/gg
	github.com/gdygd/goglib
	github.com/gdygd/goshm

test:
	go test -v

build:
	go build -o ../bin/apisvr apisvr.go

run:
	go run osvmssvr.go