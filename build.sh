#!/bin/bash
wget https://raw.githubusercontent.com/aura-studio/dynamic-cli/master/Dockerfile -O /tmp/dynamic-cli.Dockerfile
docker build --no-cache -t dynamic-cli:latest -f /tmp/dynamic-cli.Dockerfile --build-arg VERSION=latest .