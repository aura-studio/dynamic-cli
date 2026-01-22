package build

func init() {
	templateMap["{{.Dir}}/builder.sh"] = templateBuilderSh
}

const templateBuilderSh = `#!/bin/bash
set -e
mkdir -p {{.Dir}}
cd {{.Dir}}
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go clean --modcache
go mod tidy

# Resolve source module commit id (best-effort) for meta file.
# Prefer Go's module download metadata (Origin.Hash). Fallback to parsing pseudo-version.
commit_id=""
download_json=$(go mod download -json "{{.Module}}@{{.Version}}" 2>/dev/null || true)
if [ -n "$download_json" ]; then
	commit_id=$(printf "%s" "$download_json" | tr -d '\r' | sed -n 's/.*"Hash"[[:space:]]*:[[:space:]]*"\([0-9a-f]\{7,\}\)".*/\1/p' | head -n1)
fi
if [ -z "$commit_id" ]; then
	resolved_ver=$(go list -m -f '{{"{{"}}.Version{{"}}"}}' "{{.Module}}" 2>/dev/null || true)
	commit_id=$(printf "%s" "$resolved_ver" | sed -n 's/.*-\([0-9a-f]\{7,\}\)$/\1/p')
fi
if [ -z "$commit_id" ]; then
	commit_id="unknown"
fi

{{if eq .Variant "generic"}}
go build -trimpath -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN" {{.Dir}}/libcgo_{{.Name}}
go build -trimpath -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN" {{.Dir}}/libgo_{{.Name}}
{{else}}
# Unsupported variant: {{.Variant}}. Falling back to generic.
go build -trimpath -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN" {{.Dir}}/libcgo_{{.Name}}
go build -trimpath -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN" {{.Dir}}/libgo_{{.Name}}
{{end}}

# Timestamp suffix for backup artifacts and meta file.
# Format example: 20240612T153000Z
ts=$(date -u "+%Y%m%dT%H%M%SZ")

# Backup .so files with timestamp suffix.
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$ts
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$ts

# Generate meta file for this libcgo/libgo pair.
meta={{.Dir}}/meta_{{.Name}}.json
cat > "$meta" <<EOF
{
  "repo": "{{.Module}}",
  "commit_id": "$commit_id",
  "build_time": "$ts"
}
EOF

# Backup meta file with the same timestamp suffix.
cp -rf "$meta" "$meta.$ts"
`
