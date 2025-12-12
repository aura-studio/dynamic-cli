#!/usr/bin/env bash
set -euo pipefail

VERSION=$(git rev-parse --short HEAD)
IMAGE_NAME="dynamic-cli"

usage() {
	cat <<EOF
Usage: $(basename "$0") [install|docker|release|help]

Commands:
	install   Install the CLI at current git revision
	docker    Build Docker image with tag :latest and VERSION arg
	release   Build multi-platform binaries into dist/
	help      Show this help
EOF
}

install() {
	echo "Installing github.com/aura-studio/dynamic-cli@${VERSION}"
	go install "github.com/aura-studio/dynamic-cli@${VERSION}"
}

docker_build() {
	echo "Building Docker image ${IMAGE_NAME}:latest (VERSION=${VERSION})"
	docker build -t "${IMAGE_NAME}:latest" --build-arg VERSION="${VERSION}" .
}

release() {
	echo "Building release binaries (VERSION=${VERSION})"
	mkdir -p dist
	platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")
	for p in "${platforms[@]}"; do
		IFS="/" read -r GOOS GOARCH <<<"$p"
		outfile="dynamic-cli-${GOOS}-${GOARCH}"
		[[ "$GOOS" == "windows" ]] && outfile+=".exe"
		echo "  -> ${outfile}"
		GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o "dist/${outfile}" ./
	done
	echo "Artifacts in dist/"
}

cmd="${1:-install}"
case "$cmd" in
	install)
		install ;;
	docker)
		docker_build ;;
	release)
		release ;;
	help|-h|--help)
		usage ;;
	*)
		echo "Unknown command: $cmd" >&2
		usage
		exit 1 ;;
esac