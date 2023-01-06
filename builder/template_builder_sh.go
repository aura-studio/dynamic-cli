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
`
