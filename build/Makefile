MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

.EXPORT_ALL_VARIABLES:
OUT_DIR ?= _output/
BIN_DIR := $(OUT_DIR)/bin
PRJ_SRC_PATH := github.com/infinivision/citus-operator

.PHONY: all

all: citus-operator

citus-operator:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i -v -o ${BIN_DIR}/$@  ../cmd/manager/main.go

.PHONY: clean gen image

clean:
	rm -rf ${BIN_DIR}/*

gen:
	cd .. && operator-sdk generate k8s

image:
	cd .. && docker rmi citus-operator &&  operator-sdk build citus-operator
