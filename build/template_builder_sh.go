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

if [ "$actual_os" != "$expected_os" ] || [ "$actual_arch" != "$expected_arch" ] || [ "$actual_compiler" != "$expected_compiler" ]; then
	echo "Environment mismatch" >&2
	echo "  expected: OS=$expected_os ARCH=$expected_arch COMPILER=$expected_compiler" >&2
	echo "  actual:   OS=$actual_os ARCH=$actual_arch COMPILER=$actual_compiler" >&2
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
# Friendly timestamp with timezone, ISO-like
ts=$(date "+%Y-%m-%dT%H:%M:%S%z")
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$ts
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$ts
`
