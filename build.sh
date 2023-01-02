#!/bin/bash
if [ -z "$1" ]; then
	go install github.com/aura-studio/dynamic-cli@`git rev-parse --short HEAD`
elif [ "$1" == "docker" ]; then
	docker build -t dynamic-cli:latest --build-arg VERSION=`git rev-parse --short HEAD` .
else
	echo "Usage: build.sh [docker]"
fi