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

# Resolve source module commit id (best-effort) for artifact suffix naming.
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
# Trim to keep filenames short.
commit_id_short="$commit_id"
if [ "$commit_id" != "unknown" ]; then
	commit_id_short="${commit_id:0:12}"
fi
{{if eq .Variant "generic"}}
go build -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="" {{.Dir}}/libcgo_{{.Name}}
go build -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-r {{.Dir}}/" {{.Dir}}/libgo_{{.Name}}
{{else}}
# Unsupported variant: {{.Variant}}. Falling back to generic.
go build -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="" {{.Dir}}/libcgo_{{.Name}}
go build -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-r {{.Dir}}/" {{.Dir}}/libgo_{{.Name}}
{{end}}
# Friendly timestamp with timezone, ISO-like, force +08:00
# Use Asia/Shanghai timezone and insert colon into %z for +08:00
ts=$(TZ='Asia/Shanghai' date "+%Y-%m-%dT%H:%M:%S%z")
# Convert +0800 -> +08:00 for readability (portable bash string ops)
ts="${ts%??}:${ts: -2}"
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$commit_id_short.$ts
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$commit_id_short.$ts
`
