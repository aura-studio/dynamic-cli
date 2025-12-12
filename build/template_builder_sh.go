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
go build -o {{.Dir}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="" {{.Dir}}/libcgo_{{.Name}}
go build -o {{.Dir}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="-r {{.Dir}}/" {{.Dir}}/libgo_{{.Name}}
cp -rf {{.Dir}}/libgo_{{.Name}}.so {{.Dir}}/libgo_{{.Name}}.so.$(date "+%Y%m%d%H%M%S")
cp -rf {{.Dir}}/libcgo_{{.Name}}.so {{.Dir}}/libcgo_{{.Name}}.so.$(date "+%Y%m%d%H%M%S")
`
