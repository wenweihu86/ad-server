#!/bin/bash

# 下面这2行是为了下载编译依赖，只有首次编译或者依赖更新时才需要执行。
# cp -f mirrors.yaml ~/.glide/
# glide install

go build -o ad_server main/main.go
