#@IgnoreInspection BashAddShebang
export LDFLAGS="-w -s"

EXEC_NAME=gateway
DIST_DIR=bin

all: build test

build:
	mkdir -p ${DIST_DIR} && rm -f ${DIST_DIR}/* && go build -o ${DIST_DIR}/${EXEC_NAME} -race ./src/
run:
	go run -race ./src/

test:
	go test -v -race ./src/...

.PHONY: build run test container
