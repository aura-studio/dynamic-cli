package build

func init() {
	templateMap["{{.Dir}}/builder.sh"] = templateBuilderSh
}

const templateBuilderSh = `#!/bin/sh
mkdir -p {{.Dir}}
cd {{.Dir}}
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
expected_os="{{.OS}}"
expected_arch="{{.Arch}}"
expected_compiler="{{.Compiler}}"

# Detect actual environment via go env
actual_os=$(go env GOOS 2>/dev/null)
actual_arch=$(go env GOARCH 2>/dev/null)
actual_compiler=$(go env GOVERSION 2>/dev/null)

# Check toolchain syntax support by GOVERSION (Go >= 1.21)
toolchain_supported="no"
ver="${actual_compiler#go}"
major="${ver%%.*}"
minor_part="${ver#*.}"
minor="${minor_part%%.*}"
[ -z "$major" ] && major=0
[ -z "$minor" ] && minor=0
case "$major" in
	''|*[!0-9]*) major=0 ;;
esac
case "$minor" in
	''|*[!0-9]*) minor=0 ;;
esac
if [ "$major" -gt 1 ] || { [ "$major" -eq 1 ] && [ "$minor" -ge 21 ]; }; then
	toolchain_supported="yes"
fi

if [ "$actual_os" != "$expected_os" ] || [ "$actual_arch" != "$expected_arch" ] || [ "$toolchain_supported" != "yes" ]; then
	echo "Environment mismatch" >&2
	echo "  expected: OS=$expected_os ARCH=$expected_arch TOOLCHAIN=required" >&2
	echo "  actual:   OS=$actual_os ARCH=$actual_arch TOOLCHAIN=$toolchain_supported (compiler=$actual_compiler)" >&2
	exit 1
fi
go clean --modcache
go mod tidy
{{if eq .Variant "plain"}}
go build -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="" {{.Dir}}/libcgo_{{.Name}}
go build -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-r {{.Dir}}/" {{.Dir}}/libgo_{{.Name}}
{{else}}
# Unsupported variant: {{.Variant}}. Falling back to plain.
go build -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="" {{.Dir}}/libcgo_{{.Name}}
go build -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-r {{.Dir}}/" {{.Dir}}/libgo_{{.Name}}
{{end}}
# Friendly timestamp with timezone, ISO-like, force +08:00
# Use Asia/Shanghai timezone and insert colon into %z for +08:00
ts=$(TZ='Asia/Shanghai' date "+%Y-%m-%dT%H:%M:%S%z")
# Convert +0800 -> +08:00 for readability (portable bash string ops)
ts="${ts%??}:${ts: -2}"
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$ts
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$ts
`
