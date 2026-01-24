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

# Module path for meta injection
meta_module_path="{{.Module}}"

# Resolve source module commit id (best-effort) for meta file.
# Prefer Go's module download metadata (Origin.Hash). Fallback to parsing pseudo-version.
meta_commit_id=""
download_json=$(go mod download -json "$meta_module_path@{{.Version}}" 2>/dev/null || true)
if [ -n "$download_json" ]; then
	meta_commit_id=$(printf "%s" "$download_json" | tr -d '\r' | sed -n 's/.*"Hash"[[:space:]]*:[[:space:]]*"\([0-9a-f]\{7,\}\)".*/\1/p' | head -n1)
fi
if [ -z "$meta_commit_id" ]; then
	resolved_ver=$(go list -m -f '{{"{{"}}.Version{{"}}"}}' "$meta_module_path" 2>/dev/null || true)
	meta_commit_id=$(printf "%s" "$resolved_ver" | sed -n 's/.*-\([0-9a-f]\{7,\}\)$/\1/p')
fi
if [ -z "$meta_commit_id" ]; then
	meta_commit_id="unknown"
fi

# Build timestamp for meta injection (e.g. "2024-06-12 15:30:00 CST+0800")
meta_build_ts=$(TZ='Asia/Shanghai' date "+%Y-%m-%d %H:%M:%S CST%z")

# ldflags for injecting meta constants
cgo_ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN"
go_ldflags="-s -w -extldflags=-Wl,-rpath,\$ORIGIN -X 'main.MetaModulePath=$meta_module_path' -X 'main.MetaCommitID=$meta_commit_id' -X 'main.MetaBuildTs=$meta_build_ts'"

{{if eq .Variant "generic"}}
go build -trimpath -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="$cgo_ldflags" {{.Dir}}/libcgo_{{.Name}}
go build -trimpath -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="$go_ldflags" {{.Dir}}/libgo_{{.Name}}
{{else}}
# Unsupported variant: {{.Variant}}. Falling back to generic.
go build -trimpath -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="$cgo_ldflags" {{.Dir}}/libcgo_{{.Name}}
go build -trimpath -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="$go_ldflags" {{.Dir}}/libgo_{{.Name}}
{{end}}

# Backup .so files with timestamp suffix (reuse meta_build_ts).
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$meta_build_ts
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$meta_build_ts

# Generate meta file for this libcgo/libgo pair (for backward compatibility).
meta={{.Dir}}/meta_{{.Name}}.json
cat > "$meta" <<EOF
{
  "module_path": "$meta_module_path",
  "commit_id": "$meta_commit_id",
  "build_ts": "$meta_build_ts"
}
EOF

# Backup meta file with the same timestamp suffix.
cp -rf "$meta" "$meta.$meta_build_ts"
`
