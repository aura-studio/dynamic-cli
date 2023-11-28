package builder

func init() {
	templateMap["{{.House}}/{{.Name}}/builder.sh"] = builderBash
}

const builderBash = `#!/bin/sh
cd {{.House}}/{{.Name}}
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go clean --modcache
go mod tidy
go build -o {{.House}}/{{.Name}}/libcgo_{{.Name}}.so -buildvcs=false -buildmode=c-shared -ldflags="{{.LDFlags}}" {{.House}}/{{.Name}}/libcgo_{{.Name}}
go build -o {{.House}}/{{.Name}}/libgo_{{.Name}}.so -buildvcs=false -buildmode=plugin -ldflags="{{.LDFlags}}-r {{.House}}/{{.Name}}/" {{.House}}/{{.Name}}/libgo_{{.Name}}
cp -rf {{.House}}/{{.Name}}/libgo_{{.Name}}.so {{.House}}/{{.Name}}/libgo_{{.Name}}.$(date "+%Y%m%d%H%M%S")-{{.Version}}.so
cp -rf {{.House}}/{{.Name}}/libcgo_{{.Name}}.so {{.House}}/{{.Name}}/libcgo_{{.Name}}.$(date "+%Y%m%d%H%M%S")-{{.Version}}.so
`
