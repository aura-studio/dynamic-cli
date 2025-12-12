package build

func init() {
	templateMap["{{.Dir}}/builder.sh"] = builderBash
}

const builderBash = `#!/bin/sh
cd {{.Dir}}
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go clean --modcache
go mod tidy
{{/* Variant-based build commands; currently support plain */}}
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
