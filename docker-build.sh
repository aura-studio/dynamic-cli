#!/bin/bash
docker build -t dynamic-cli:latest --build-arg VERSION=`git rev-parse --short HEAD` .