package build

func init() {
	templateMap["{{.Dir}}/builder.sh"] = templateBuilderSh
}

const templateBuilderSh = `#!/bin/bash
mkdir -p {{.Dir}}
cd {{.Dir}}
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go clean --modcache
go mod tidy
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
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$ts
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$ts
`
