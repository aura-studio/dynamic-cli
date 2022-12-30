package builder

func init() {
	templateMap["{{.House}}/{{.Name}}_{{.Version}}/builder.sh"] = builderBash
}

const builderBash = `#!/bin/sh
export CGO_ENABLED=1
export GO111MODULE=on
export GOPRIVATE={{.Module}}
go mod tidy
go build -o {{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}}.so -buildvcs=false -buildmode=c-shared -ldflags="{{.LDFlags}}" {{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}} && upx -1 {{.House}}/{{.Name}}_{{.Version}}/libcgo_{{.Name}}_{{.Version}}.so
go build -o {{.House}}/{{.Name}}_{{.Version}}/libgo_{{.Name}}_{{.Version}}.so -buildvcs=false -buildmode=plugin -ldflags="{{.LDFlags}}-r {{.House}}/{{.Name}}_{{.Version}}/" {{.House}}/{{.Name}}_{{.Version}}/libgo_{{.Name}}_{{.Version}} && upx -1 {{.House}}/{{.Name}}_{{.Version}}/libgo_{{.Name}}_{{.Version}}.so
`
