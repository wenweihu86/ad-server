#!/bin/bash

glide install
go build -o ad_server src/main/main.go
