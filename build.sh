#!/bin/bash

export GOPATH=${GOPATH}:`pwd`
echo $GOPATH
go build -o ad_server src/main/main.go
