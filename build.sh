#!/bin/bash
docker build --no-cache -t dynamic-cli:latest --build-arg VERSION=stable .